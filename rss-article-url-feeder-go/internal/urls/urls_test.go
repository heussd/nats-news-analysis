package urls

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSize(t *testing.T) {
	assert.Equal(t, 3, len(Urls))
}

func TestIterate(t *testing.T) {
	for i, v := range Urls {
		fmt.Printf("%d %s\n", i, v)
	}
}
