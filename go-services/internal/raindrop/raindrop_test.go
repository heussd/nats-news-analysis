package raindrop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	res, err := Add("https://www.tagesschau.de/wirtschaft/konjunktur/inflationsrate-juli-verbraucherpreise-100.html")
	assert.NoError(t, err)
	assert.Equal(t, true, res)

}
