package nats

import (
	"fmt"
	"github.com/heussd/nats-raindrop-integration.go/internal/config"
	"github.com/nats-io/nats.go"
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

}

func WithMatchUrls(f func(m *nats.Msg)) {
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
