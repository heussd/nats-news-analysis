package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLocalMakeIT(t *testing.T) {
	WithArticleUrls(func(m *nats.Msg) {
		data := string(m.Data)
		fmt.Printf("Received from %s\n", data)
		assert.Equal(t, "https://www.tagesschau.de/", data)
	})
	PushToPocket("HELLO WORL2D2")
}
