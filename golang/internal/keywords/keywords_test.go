package keywords

import (
	"testing"

	"github.com/heussd/nats-news-analysis/internal/model"
	"github.com/stretchr/testify/assert"
)

func hasMatch(matched []model.Keyword, err error) bool {
	return len(matched) > 0 && err == nil
}
func firstMatch(matched []model.Keyword, err error) string {
	if hasMatch(matched, err) {
		return matched[0].Text
	}
	return ""
}

func TestLocalIT(t *testing.T) {
	assert.Equal(t, true, hasMatch(Match("Peach")))
	assert.Equal(t, false, hasMatch(Match("Pineapple")))
	assert.Equal(t, false, hasMatch(Match("Hamburger")))
	assert.Equal(t, true, hasMatch(Match("Apple")))
	assert.Equal(t, true, hasMatch(Match("Banana split")))
	assert.Equal(t, true, hasMatch(Match("Delicious pineapple recipes")))
	assert.Equal(t, true, hasMatch(Match("Delicious recipes")))
	assert.Equal(t, true, hasMatch(Match("Delicious pineapple pies")))
	assert.Equal(t, true, hasMatch(Match("Delicious dark-chocolate pineapple pies")))
	assert.Equal(t, true, hasMatch(Match("ICE cold cream")))
	assert.Equal(t, false, hasMatch(Match("ICE and also some cream")))
	assert.Equal(t, false, hasMatch(Match("whipped cream")))
	assert.Equal(t, false, hasMatch(Match("# Should not match")))

	assert.Equal(t, false, hasMatch(Match("Mister Cool")))
	assert.Equal(t, false, hasMatch(Match("Miss Gray")))
	assert.Equal(t, true, hasMatch(Match("Mississippi")))

	assert.Equal(t, false, hasMatch(Match("Bias")))
	assert.Equal(t, true, hasMatch(Match("as")))

	assert.Equal(t, true, hasMatch(Match("All of us")))
	assert.Equal(t, true, hasMatch(Match("All-of-us")))
	assert.Equal(t, false, hasMatch(Match("Alloofuus")))

	assert.Equal(t, false, hasMatch(Match("I drink cold beer. I eat hot pizza.")))
	assert.Equal(t, true, hasMatch(Match("I ate cold yummy pizza yesterday afternoon. I drink hot chocolate.")))

	assert.Equal(t, true, hasMatch(Match("The king lived long and prosper.")))
	assert.Equal(t, true, hasMatch(Match("Long live the king.")))
	assert.Equal(t, true, hasMatch(Match("The queen lived long and prosper.")))
	assert.Equal(t, true, hasMatch(Match("Long live the queen.")))

	assert.Equal(t, false, hasMatch(Match("Like king and queen.")))

}

func TestStringMatchReturn(t *testing.T) {
	text := firstMatch(Match("A little Peach a day"))
	assert.Equal(t, "(?i)\\b(Apple|peach)", text)

	text = firstMatch(Match("I like to eat delicious original organic-sourced pineapple pies twice a day"))
	assert.Equal(t, "(?i)(delicious).*(pie|recipes)", text)

	text = firstMatch(Match("Long live the queen. Something else"))
	assert.Equal(t, "(?i)^(?=.*(king|queen))(?=.*long).*", text)

}
