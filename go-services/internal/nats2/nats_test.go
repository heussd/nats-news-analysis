package nats

import (
	"fmt"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestPush(t *testing.T) {
	Publish("https://www.tagesschau.de/2")
}

func TestKV(t *testing.T) {
	key := "OMG21223232"

	assert.Equal(t, false, hasKV(key))
	putKV(key)
	assert.Equal(t, true, hasKV(key))
}

func TestKV2(t *testing.T) {
	key := "https://www.tagesschau.de/inland/gesellschaft/food-sharing-101.html"

	assert.Equal(t, false, hasKV(key))
	putKV(key)
	assert.Equal(t, true, hasKV(key))
}

func TestPull(t *testing.T) {
	WithFeedUrls(func(m *nats.Msg) {
		data := string(m.Data)
		fmt.Printf("Received from %s\n", data)
		assert.Equal(t, "https://www.tagesschau.de/", data)
	})
}
