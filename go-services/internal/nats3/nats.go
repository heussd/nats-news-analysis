package nats

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"

	"nats-news-analysis/internal/config"
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

}
