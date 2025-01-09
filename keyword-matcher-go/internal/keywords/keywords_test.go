package keywords

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, false, first(Match("I drink cold beer. I eat hot pizza.")))
	assert.Equal(t, true, first(Match("I ate cold yummy pizza yesterday afternoon. I drink hot chocolate.")))

	assert.Equal(t, true, first(Match("The king lived long and prosper.")))
	assert.Equal(t, true, first(Match("Long live the king.")))
	assert.Equal(t, true, first(Match("The queen lived long and prosper.")))
	assert.Equal(t, true, first(Match("Long live the queen.")))

	assert.Equal(t, false, first(Match("Like king and queen.")))

}

func TestStringMatchReturn(t *testing.T) {
	_, text := Match("A little Peach a day")
	assert.Equal(t, "Apple peach", text)

	_, text = Match("I like to eat delicious original organic-sourced pineapple pies twice a day")
	assert.Equal(t, "delicious pie recipes", text)

	_, text = Match("Long live the queen. Something else")
	assert.Equal(t, "king queen long", text)

}
