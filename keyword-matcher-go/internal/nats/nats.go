package nats

import (
	"github.com/heussd/nats-news-keyword-matcher.go/internal/config"
	"github.com/nats-io/nats.go"
	"sync"
	"time"
)

var nc, _ = nats.Connect(config.NatsServer)
var js, _ = nc.JetStream(nats.PublishAsyncMaxPending(256))

func init() {
	js.AddStream(&nats.StreamConfig{
		Name:       config.NatsOutputQueueName,
		Subjects:   []string{config.NatsOutputQueueSubject},
		Retention:  nats.WorkQueuePolicy,
		Duplicates: time.Hour * 24 * 30,
	})
}

func WithArticleUrls(f func(m *nats.Msg)) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	js.Subscribe(config.NatsInputQueueSubject, f)
	wg.Done()
	wg.Wait()
}

func PushToPocket(s string) {
	headers := nats.Header{}
	headers.Add(nats.MsgIdHdr, "news-keyword-matcher"+s)

	msg := &nats.Msg{
		Header:  headers,
		Subject: config.NatsOutputQueueSubject,
		Data:    []byte(s),
	}

	js.PublishMsg(msg)
}
