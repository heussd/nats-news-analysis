package nats

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
