package keywords

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLocalIT(t *testing.T) {
	assert.Equal(t, true, Match("Peach"))
	assert.Equal(t, false, Match("Pineapple"))
	assert.Equal(t, false, Match("Hamburger"))
	assert.Equal(t, true, Match("Apple"))
	assert.Equal(t, true, Match("Banana split"))
	assert.Equal(t, true, Match("Delicious pineapple recipes"))
	assert.Equal(t, true, Match("Delicious recipes"))
	assert.Equal(t, true, Match("Delicious pineapple pies"))
	assert.Equal(t, true, Match("Delicious dark-chocolate pineapple pies"))
	assert.Equal(t, true, Match("ICE cold cream"))
	assert.Equal(t, false, Match("ICE and also some cream"))
	assert.Equal(t, false, Match("whipped cream"))
	assert.Equal(t, false, Match("# Should not match"))

	assert.Equal(t, false, Match("Mister Cool"))
	assert.Equal(t, false, Match("Miss Gray"))
	assert.Equal(t, true, Match("Mississippi"))

	assert.Equal(t, false, Match("Bias"))
	assert.Equal(t, true, Match("as"))

	assert.Equal(t, true, Match("All of us"))
	assert.Equal(t, true, Match("All-of-us"))
	assert.Equal(t, false, Match("Alloofuus"))
}
