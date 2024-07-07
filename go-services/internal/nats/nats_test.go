package nats

import (
	"fmt"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"

	"nats-news-analysis/internal/model"
)

func TestPull(t *testing.T) {
	WithArticleUrls(func(m *nats.Msg) {
		data := string(m.Data)
		fmt.Printf("Received from %s\n", data)
		assert.Equal(t, "https://www.tagesschau.de/", data)
	})
}

func TestPush(t *testing.T) {
	PushToPocket(model.Match{Url: "https://www.tagesschau.de/2", RegexId: "TEXT"})
}
