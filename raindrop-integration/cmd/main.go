package main

import (
	"encoding/json"
	"fmt"
	queue "github.com/heussd/nats-raindrop-integration.go/internal/nats"
	"github.com/heussd/nats-raindrop-integration.go/internal/raindrop"
	"github.com/nats-io/nats.go"
	"os"
	"os/signal"
	"syscall"
)

type Match struct {
	RegexId string
	Url     string
}

func main() {

	queue.WithMatchUrls(func(m *nats.Msg) {
		var match Match

		if err := json.Unmarshal([]byte(string(m.Data)), &match); err != nil {
			fmt.Errorf("failed to unmarshall %w", err)
		}

		if ok, _ := raindrop.Add(match.Url); !ok {
			fmt.Errorf("received error from raindrop\n")
		} else {
			fmt.Printf("Added to Raindrop: %s\n	", match.Url)
		}
	})

	// https://callistaenterprise.se/blogg/teknik/2019/10/05/go-worker-cancellation/
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan // Blocks here until either SIGINT or SIGTERM is received.
}
