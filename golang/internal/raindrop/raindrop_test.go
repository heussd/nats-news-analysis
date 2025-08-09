package raindrop

import (
	"testing"

	"github.com/heussd/nats-news-analysis/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	m := model.Match{
		Url: "https://www.tagesschau.de/wirtschaft/konjunktur/inflationsrate-juli-verbraucherpreise-100.html",
		Keywords: []model.Keyword{
			{
				Text: "(?i)^(?=.*(king|queen))(?=.*long).* as regex",
			},
			{
				Text: "(?i)(delicious).*(pie|recipes)",
			}},
	}
	err := Add(&m)
	assert.NoError(t, err)
}
