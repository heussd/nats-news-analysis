package nats

import (
	"github.com/heussd/nats-news-keyword-matcher.go/internal/utils"
	"github.com/nats-io/nats.go"
	"time"
)

var natsServer = utils.GetEnv("NATS_SERVER", nats.DefaultURL)
var natsInputQueueName = utils.GetEnv("NATS_INPUT_QUEUE_NAME", "article_urls")
var natsInputQueueSubject = utils.GetEnv("NATS_INPUT_QUEUE_SUBJECT", "article_url")
var natsOutputQueueName = utils.GetEnv("NATS_OUTPUT_QUEUE_NAME", "match_urls")
var natsOutputQueueSubject = utils.GetEnv("NATS_OUTPUT_QUEUE_SUBJECT", "match_url")

var nc, _ = nats.Connect(natsServer)
var js, _ = nc.JetStream(nats.PublishAsyncMaxPending(256))

func init() {
	js.AddStream(&nats.StreamConfig{
		Name:       natsOutputQueueName,
		Subjects:   []string{natsOutputQueueSubject},
		Retention:  nats.WorkQueuePolicy,
		Duplicates: time.Hour * 24 * 30,
	})
}

func WithArticleUrls(f func(m *nats.Msg)) {
	js.Subscribe(natsInputQueueSubject, f)
}

func PushToPocket(s string) {

	headers := nats.Header{}
	headers.Add(nats.MsgIdHdr, "news-keyword-matcher"+s)

	msg := &nats.Msg{
		Header:  headers,
		Subject: natsOutputQueueSubject,
		Data:    []byte(s),
	}

	js.PublishMsg(msg)
}
