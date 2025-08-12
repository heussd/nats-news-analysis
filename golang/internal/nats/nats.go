package nats

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/heussd/nats-news-analysis/internal/model"
	"github.com/heussd/nats-news-analysis/pkg/utils"
	"github.com/nats-io/nats.go"
)

var (
	nc *nats.Conn
	js nats.JetStreamContext
	kv nats.KeyValue
)

var (
	NatsServer                      = utils.GetEnv("NATS_URL", nats.DefaultURL)
	NatsPullConsumerBatchSize       = utils.GetEnv("NATS_CONSUMER_BATCH_SIZE", "5")
	NatsPullConsumerBatchSizeInt, _ = strconv.Atoi(NatsPullConsumerBatchSize)
	NatsKeyValueBucket              = utils.GetEnv("NATS_KV_BUCKET", "article-urls-proposed")
	DefaultDupeWindow               = time.Hour * 24 * 30
)

func init() {
	for nc == nil {
		var err error
		nc, err = nats.Connect(NatsServer)
		if err != nil {
			fmt.Printf("Could not connect to NATS server at %s: %s\n", NatsServer, err)
			time.Sleep(5 * time.Second)
		}
	}

	var err error
	if js, err = nc.JetStream(nats.PublishAsyncMaxPending(256)); err != nil {
		panic(err)
	}

	if kv, err = js.KeyValue(NatsKeyValueBucket); err != nil {
		if kv, err = js.CreateKeyValue(&nats.KeyValueConfig{
			Bucket: NatsKeyValueBucket,
			TTL:    time.Hour * 24 * 365,
		}); err != nil {
			panic(err)
		}
	}
}

func Subscribe[T model.PayloadTypes](
	stream *nats.StreamInfo,
	consumer *nats.ConsumerInfo,
	f func(m *T),
	loopForever bool,
) (err error) {
	var sub *nats.Subscription
	if sub, err = js.PullSubscribe(stream.Config.Name, consumer.Config.Durable); err != nil {
		return fmt.Errorf("error doing request: %w", err)
	}

	for {
		msgs, err := sub.Fetch(
			NatsPullConsumerBatchSizeInt,
			nats.MaxWait(60*time.Second))
		if err != nil {
			if err == nats.ErrTimeout {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			log.Printf("Fetch error: %v", err)
			panic(err)
		}

		for _, msg := range msgs {
			msg.Ack()

			var obj *T

			if err := json.Unmarshal([]byte(string(msg.Data)), &obj); err != nil {
				return fmt.Errorf("failed to unmarshall %w", err)
			}
			f(obj)
		}

		if !loopForever {
			break
		}
	}

	return nil
}

func Publish[T model.PayloadTypes](
	stream *nats.StreamInfo,
	message T,
	natsMessageId string, // Meaningful ID to make use of NATS's deduplication feature
	persistDeduplication bool,
) (*nats.PubAck, error) {
	// Workaround for https://github.com/nats-io/nats-server/issues/3272
	// Use KV storage for remembering what we already put on the queue.
	if persistDeduplication && hasKV(message.GetUrl()) {
		return nil, nil
	}

	headers := nats.Header{}

	if natsMessageId != "" {
		headers.Add(nats.MsgIdHdr, natsMessageId)
	}

	var data []byte
	var err error
	if data, err = json.Marshal(message); err != nil {
		return nil, fmt.Errorf("cannot marshal: %w", err)
	}

	msg := &nats.Msg{
		Header:  headers,
		Subject: stream.Config.Name,
		Data:    data,
	}

	var pubAck *nats.PubAck
	if pubAck, err = js.PublishMsg(msg); err != nil {
		return nil, fmt.Errorf("publish failed: %w", err)
	}

	if persistDeduplication && err == nil {
		putKV(message.GetUrl())
	}

	return pubAck, nil
}

func AddStream(name string, dupe_window time.Duration) (str *nats.StreamInfo, err error) {
	if str, err = js.AddStream(&nats.StreamConfig{
		Name: name,
		Subjects: []string{
			name,
		},
		Retention:  nats.LimitsPolicy,
		MaxAge:     time.Hour * 24 * 90,
		Duplicates: dupe_window,
	}); err != nil {
		if err == nats.ErrStreamNameAlreadyInUse {
			fmt.Printf("Stream %s already exists with different config\n", name)
			return nil, err
		}
		return nil, err
	}
	fmt.Printf("Stream %s configured successfully\n", name)
	return str, nil
}

func AddConsumer(streamName, consumerName string) (consumer *nats.ConsumerInfo, err error) {
	if consumer, err = js.AddConsumer(
		streamName,
		&nats.ConsumerConfig{
			Durable:       consumerName,
			DeliverPolicy: nats.DeliverNewPolicy,
			AckPolicy:     nats.AckExplicitPolicy,
		}); err != nil {
		if err == nats.ErrConsumerNameAlreadyInUse {
			fmt.Printf("Consumer %s already exists with different config\n", consumerName)
			return nil, err
		}
		return nil, err
	}
	fmt.Printf("Consumer %s configured successfully for stream %s\n", consumerName, streamName)

	return consumer, nil
}

func AddStreamOrDie(streamName string, dupe_window time.Duration) (stream *nats.StreamInfo) {
	var err error
	if stream, err = AddStream(streamName, dupe_window); err != nil {
		panic(fmt.Errorf("failed to add stream %s: %w", streamName, err))
	}
	return stream
}

func AddConsumerOrDie(stream *nats.StreamInfo, consumerName string) (consumer *nats.ConsumerInfo) {
	var err error
	if consumer, err = AddConsumer(stream.Config.Name, consumerName); err != nil {
		panic(fmt.Errorf("failed to add consumer %s for stream %s: %v", consumerName, stream.Config.Name, err))
	}
	return consumer
}

func generateMessageId(prefix string, match model.Match) string {
	hash := sha1.New()
	io.WriteString(hash, match.Url)

	for _, entry := range match.Keywords {
		io.WriteString(hash, entry.Text)
	}
	hexHash := hex.EncodeToString(hash.Sum(nil))

	return fmt.Sprintf("%s-%s", prefix, hexHash)
}

func workaroundNatsGoKVClientKeyBug(key string) string {
	hash := sha1.New()
	io.WriteString(hash, key)

	return hex.EncodeToString(hash.Sum(nil))
}

func putKV(key string) {
	kv.Put(workaroundNatsGoKVClientKeyBug(key), []byte("1"))
}

func hasKV(key string) bool {
	if v, err := kv.Get(workaroundNatsGoKVClientKeyBug(key)); err == nil {
		value := string(v.Value())
		if value == "1" {
			return true
		}
	}
	return false
}
