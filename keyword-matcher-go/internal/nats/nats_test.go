package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPull(t *testing.T) {
	WithArticleUrls(func(m *nats.Msg) {
		data := string(m.Data)
		fmt.Printf("Received from %s\n", data)
		assert.Equal(t, "https://www.tagesschau.de/", data)
	})
}

func TestPush(t *testing.T) {
	PushToPocket("https://www.tagesschau.de/2", "TEXT")
}
