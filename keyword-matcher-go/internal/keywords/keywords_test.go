package keywords

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func first(flag bool, _ string) bool {
	return flag
}

func TestLocalIT(t *testing.T) {
	assert.Equal(t, true, first(Match("Peach")))
	assert.Equal(t, false, first(Match("Pineapple")))
	assert.Equal(t, false, first(Match("Hamburger")))
	assert.Equal(t, true, first(Match("Apple")))
	assert.Equal(t, true, first(Match("Banana split")))
	assert.Equal(t, true, first(Match("Delicious pineapple recipes")))
	assert.Equal(t, true, first(Match("Delicious recipes")))
	assert.Equal(t, true, first(Match("Delicious pineapple pies")))
	assert.Equal(t, true, first(Match("Delicious dark-chocolate pineapple pies")))
	assert.Equal(t, true, first(Match("ICE cold cream")))
	assert.Equal(t, false, first(Match("ICE and also some cream")))
	assert.Equal(t, false, first(Match("whipped cream")))
	assert.Equal(t, false, first(Match("# Should not match")))

	assert.Equal(t, false, first(Match("Mister Cool")))
	assert.Equal(t, false, first(Match("Miss Gray")))
	assert.Equal(t, true, first(Match("Mississippi")))

	assert.Equal(t, false, first(Match("Bias")))
	assert.Equal(t, true, first(Match("as")))

	assert.Equal(t, true, first(Match("All of us")))
	assert.Equal(t, true, first(Match("All-of-us")))
	assert.Equal(t, false, first(Match("Alloofuus")))
}

func TestLocalIT2(t *testing.T) {
	_, text := Match("Peach")
	assert.Equal(t, "(?i)\\b(Apple|peach)", text)
}
