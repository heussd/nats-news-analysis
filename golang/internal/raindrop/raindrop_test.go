package raindrop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	err := Add("https://www.tagesschau.de/wirtschaft/konjunktur/inflationsrate-juli-verbraucherpreise-100.html")
	assert.NoError(t, err)
}
