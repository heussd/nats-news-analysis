package nats

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/shomali11/util/xhashes"

	"nats-news-analysis/internal/config"
	"nats-news-analysis/internal/model"
)

var nc *nats.Conn
var js nats.JetStreamContext
var kv nats.KeyValue

func init() {
	var err error

	for nc == nil {
		var err error
		nc, err = nats.Connect(config.NatsServer)
		if err != nil {
			fmt.Printf("Could not connect to NATS server at %s, because: ", config.NatsServer)
			fmt.Println(err)
			time.Sleep(5 * time.Second)
		}
	}

	if js, err = nc.JetStream(nats.PublishAsyncMaxPending(256)); err != nil {
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

func WithArticleUrls(f func(m *nats.Msg)) {
	js.AddStream(&nats.StreamConfig{
		Name:       config.NatsOutputQueueNameArticle,
		Subjects:   []string{config.NatsOutputQueueSubjectArticle},
		Retention:  nats.WorkQueuePolicy,
		Duplicates: time.Hour * 24 * 30,
	})

	wg := sync.WaitGroup{}
	wg.Add(1)

	var subscribed = false
	for !subscribed {
		_, err := js.QueueSubscribe(config.NatsInputQueueSubjectArticle, config.NatsInputQueueNameArticle, f)
		if err != nil {
			fmt.Printf("Cannot subscribe queue %s subject %s, because: ", config.NatsInputQueueNameArticle, config.NatsInputQueueSubjectArticle)
			fmt.Println(err)
			time.Sleep(5 * time.Second)
		} else {
			subscribed = true
		}
	}

	wg.Done()
	wg.Wait()
}

func PushToPocket(match model.Match) {
	headers := nats.Header{}
	headers.Add(nats.MsgIdHdr, "news-keyword-matcher"+match.Url)

	data, _ := json.Marshal(match)

	msg := &nats.Msg{
		Header:  headers,
		Subject: config.NatsOutputQueueSubjectMatch,
		Data:    data,
	}

	js.PublishMsg(msg)
}

func WithMatchUrls(f func(m *nats.Msg)) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	var subscribed = false
	for !subscribed {
		_, err := js.QueueSubscribe(config.NatsInputQueueSubjectMatch, config.NatsInputQueueNameMatch, f)
		if err != nil {
			fmt.Printf("Cannot subscribe queue %s subject %s, because: ", config.NatsInputQueueNameMatch, config.NatsInputQueueSubjectMatch)
			fmt.Println(err)
			time.Sleep(5 * time.Second)
		} else {
			subscribed = true
		}
	}

	wg.Done()
	wg.Wait()
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
func PublishMatch(url string) {
	Publish(config.NatsInputQueueSubjectMatch, url)
}
func PublishArticle(url string) {
	Publish(config.NatsInputQueueSubjectArticle, url)
}

func Publish(subject string, url string) (js nats.JetStreamContext, err error) {
	if js, err = nc.JetStream(nats.PublishAsyncMaxPending(256)); err != nil {
		panic(err)
	}

	// Workaround for https://github.com/nats-io/nats-server/issues/3272
	// Use KV storage for remembering what we already put on the queue.
	if !hasKV(url) {

		fmt.Printf("Publishing %s\n", url)

		headers := nats.Header{}
		headers.Add(nats.MsgIdHdr, url)

		msg := &nats.Msg{
			Header:  headers,
			Subject: subject,
			Data:    []byte(url),
		}

		js.PublishMsg(msg)
		putKV(url)
	}

	return js, nil
}
