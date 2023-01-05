package urls

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
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

func TestRandomness(t *testing.T) {
	ReloadUrls()
	var first = strings.Join(Urls, "-")
	fmt.Println("URLS: ", first)

	ReloadUrls()
	var second = strings.Join(Urls, "-")
	fmt.Println("URLS: ", second)

	ReloadUrls()
	var third = strings.Join(Urls, "-")
	fmt.Println("URLS: ", third)

	assert.NotEqual(t, first+second+third, third+second+first)
}
