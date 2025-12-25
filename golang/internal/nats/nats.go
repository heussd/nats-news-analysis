// Package nats provides abstractions for interacting with NATS messaging service.
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

	for key, value := range StreamConfigs {
		fmt.Printf("Init stream %s with %+v\n", key, value)
		if _, err := addStream(value); err != nil {
			panic(err)
		}
	}
}

func Subscribe[T model.PayloadTypes](
	f func(m *T),
	opts ...func(*NatsSubscribeOpts),
) (err error) {
	props := GetDefaultStreamConfig[T]()
	for _, opt := range opts {
		opt(props)
	}

	fmt.Printf("Setting up NATS consumer with %+v\n", props)

	if _, err = addConsumer(*props); err != nil {
		return fmt.Errorf("failed to add consumer: %w", err)
	}

	var sub *nats.Subscription
	if sub, err = js.PullSubscribe(props.SubjectName, props.ConsumerName); err != nil {
		return fmt.Errorf("failed to pull subscribe: %w", err)
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

		if props.TerminateAfterOneMessage {
			break
		}
	}

	return nil
}

func Publish[T model.PayloadTypes](
	message T,
	opts ...func(*NatsPublishOpts),
) (*nats.PubAck, error) {
	props := defaultNatsPublishOpts
	for _, opt := range opts {
		opt(props)
	}

	// Workaround for https://github.com/nats-io/nats-server/issues/3272
	// Use KV storage for remembering what we already put on the queue.
	if props.PersistDeduplication && hasKV(message.GetUrl()) {
		return nil, nil
	}

	headers := nats.Header{}

	if props.NatsMessageID != "" {
		headers.Add(nats.MsgIdHdr, props.NatsMessageID)
	}

	var data []byte
	var err error
	if data, err = json.Marshal(message); err != nil {
		return nil, fmt.Errorf("cannot marshal: %w", err)
	}

	msg := &nats.Msg{
		Header:  headers,
		Subject: props.Subject,
		Data:    data,
	}

	var pubAck *nats.PubAck
	if pubAck, err = js.PublishMsg(msg); err != nil {
		return nil, fmt.Errorf("publish failed: %w", err)
	}

	if props.PersistDeduplication {
		putKV(message.GetUrl())
	}

	return pubAck, nil
}

func addStream(props NatsStreamOpts) (str *nats.StreamInfo, err error) {
	cfg := &nats.StreamConfig{
		Name:       props.StreamName,
		Subjects:   []string{props.SubjectName},
		Retention:  nats.LimitsPolicy,
		MaxAge:     time.Hour * 24 * 90,
		Duplicates: props.DupeWindow,
	}
	if str, err = js.AddStream(cfg); err != nil {
		if err == nats.ErrStreamNameAlreadyInUse {
			fmt.Printf("Stream %s already exists with different config, reconfiguring...\n", props.StreamName)
			if str, err = js.UpdateStream(cfg); err != nil {
				fmt.Printf("Stream reconfiguration failed: %s", err)
				return nil, err
			}
			return str, nil
		}
		return nil, err
	}
	fmt.Printf("✅ Stream %s configured successfully\n", props.StreamName)
	return str, nil
}

func addConsumer(props NatsSubscribeOpts) (consumer *nats.ConsumerInfo, err error) {
	if consumer, err = js.AddConsumer(
		props.StreamName,
		&nats.ConsumerConfig{
			Durable:       props.ConsumerName,
			DeliverPolicy: nats.DeliverNewPolicy,
			AckPolicy:     nats.AckExplicitPolicy,
		}); err != nil {
		if err == nats.ErrConsumerNameAlreadyInUse {
			fmt.Printf("Consumer %s already exists with different config\n", props.ConsumerName)
			return nil, err
		}
		return nil, err
	}
	fmt.Printf("Consumer %s configured successfully for stream %s\n", props.ConsumerName, props.StreamName)

	return consumer, nil
}

func generateMessageId(prefix string, match model.Match) string {
	hash := sha1.New()
	_, _ = io.WriteString(hash, match.Url)

	for _, entry := range match.Keywords {
		_, _ = io.WriteString(hash, entry.Text)
	}
	hexHash := hex.EncodeToString(hash.Sum(nil))

	return fmt.Sprintf("%s-%s", prefix, hexHash)
}

func workaroundNatsGoKVClientKeyBug(key string) string {
	hash := sha1.New()
	_, _ = io.WriteString(hash, key)

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

func WaitFor(subject string) {
	var s string
	for {
		if s, _ = js.StreamNameBySubject(subject); s != "" {
			break
		}
		fmt.Printf("Waiting for a stream to accept subject \"%s\"\n", subject)
		time.Sleep(time.Second * 2)
	}
}
