package nats

import (
	"fmt"
	"github.com/heussd/rss-article-url-feeder.go/internal/config"
	"github.com/nats-io/nats.go"
	"github.com/shomali11/util/xhashes"
	"sync"
	"time"
)

var nc *nats.Conn
var js nats.JetStreamContext
var kv nats.KeyValue

func init() {
	var err error

	for nc == nil {
		if nc, err = nats.Connect(config.NatsServer); err != nil {
			fmt.Printf("Could not connect to NATS server at %s, because: ", config.NatsServer)
			fmt.Println(err)
			time.Sleep(5 * time.Second)
		}
	}

	if js, err = nc.JetStream(nats.PublishAsyncMaxPending(256)); err != nil {
		fmt.Println("HERE")
		panic(err)
	}

	if _, err = js.AddStream(&nats.StreamConfig{
		Name: config.NatsOutputQueueName,
		Subjects: []string{
			config.NatsOutputQueueSubject,
		},
		Retention:  nats.WorkQueuePolicy,
		Duplicates: time.Hour * 24 * 30,
	}); err != nil {

		panic(err)
	}

	if kv, err = js.KeyValue(config.NatsKeyValueBucket); err != nil {
		if kv, err = js.CreateKeyValue(&nats.KeyValueConfig{
			Bucket: config.NatsKeyValueBucket,
		}); err != nil {
			panic(err)
		}
	}
}

func workaroundNatsGoKVClientKeyBug(key string) string {
	return xhashes.SHA1(key)
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

func Publish(url string) {
	// Workaround for https://github.com/nats-io/nats-server/issues/3272
	// Use KV storage for remembering what we already put on the queue.
	if !hasKV(url) {

		fmt.Printf("Publishing %s\n", url)

		headers := nats.Header{}
		headers.Add(nats.MsgIdHdr, url)

		msg := &nats.Msg{
			Header:  headers,
			Subject: config.NatsOutputQueueSubject,
			Data:    []byte(url),
		}

		js.PublishMsg(msg)
		putKV(url)
	}
}

func WithFeedUrls(f func(m *nats.Msg)) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	var subscribed = false
	for !subscribed {
		_, err := js.QueueSubscribe(config.NatsInputQueueSubject, config.NatsInputQueueName, f)
		if err != nil {
			fmt.Printf("Cannot subscribe queue %s subject %s, because: ", config.NatsInputQueueName, config.NatsInputQueueSubject)
			fmt.Println(err)
			time.Sleep(5 * time.Second)
		} else {
			subscribed = true
		}
	}

	wg.Done()
	wg.Wait()
}
