package nats

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/dlclark/regexp2"
	"github.com/heussd/nats-news-analysis/internal/model"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestPushAndPull(t *testing.T) {
	stream := "test-stream3"

	s, _ := AddStream(stream)
	c, _ := AddConsumer(stream, "consumer")

	url := "https://www.tagesschau.de/"
	Publish(
		s,
		model.Match{
			Url: url,
			Keywords: []model.Keyword{
				{
					Text:   fmt.Sprintf("Test-%d", rand.Intn(1000)), // Avoid deduplication with random
					Id:     fmt.Sprintf("test-id-%d", rand.Intn(10000)),
					Regexp: *regexp2.MustCompile("(?i)^(?=.*(king|queen))(?=.*long).*", 0),
				},
			}},
		"",
		false,
	)
	Subscribe(s, c,
		func(m *model.Match) {
			fmt.Printf("%+v\n", m)
			assert.Equal(t, url, m.Url)
		},
		false,
	)
}

func TestAddStream(t *testing.T) {
	testStreamName := "test-stream"
	defer func() {
		// Clean up the test stream after the test
		if err := js.DeleteStream(testStreamName); err != nil && err != nats.ErrStreamNotFound {
			t.Fatalf("Failed to delete test stream: %v", err)
		}
	}()

	s, err := AddStream(testStreamName)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Equal(t, testStreamName, s.Config.Name)

	s, err = AddStream(testStreamName)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Equal(t, testStreamName, s.Config.Name)
}

func TestGenerateMessageId(t *testing.T) {
	match := model.Match{
		Url: "https://example.com",
		Keywords: []model.Keyword{
			{
				Text:   "Test Keyword",
				Id:     "test-id",
				Regexp: *regexp2.MustCompile("(?i)^(?=.*(king|queen))(?=.*long).*", 0),
			},
		},
	}

	id := generateMessageId("test-prefix", match)
	assert.NotEmpty(t, id)

	id4 := generateMessageId("test-prefix", match)
	assert.Equal(t, id, id4, "Message IDs should be the same for the same match")

	match2 := match
	match2.Keywords[0].Text = "Something else"
	id2 := generateMessageId("test-prefix", match2)
	assert.NotEqual(t, id, id2, "Message IDs should be different for different keyword texts")

	match3 := match
	match3.Url = "https://another-example.com"
	id3 := generateMessageId("test-prefix", match3)
	assert.NotEqual(t, id, id3, "Message IDs should be different for different URLs")

}
