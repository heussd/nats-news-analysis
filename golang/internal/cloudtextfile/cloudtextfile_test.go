package cloudtextfile

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var cloudTextFileUrl = "https://raw.githubusercontent.com/heussd/nats-news-analysis/refs/heads/main/golang/internal/keywords/keywords.txt"

func TestCloudTextFile(t *testing.T) {
	start := time.Now()
	lines, fromCache, err := CachedCloudTextFile(cloudTextFileUrl)
	firstRetrievalTime := time.Since(start)
	t.Logf("CachedCloudTextFile took %s", firstRetrievalTime)

	assert.NoError(t, err)
	assert.False(t, fromCache)
	assert.Greater(t, len(lines), 0)

	start = time.Now()
	lines, fromCache, err = CachedCloudTextFile(cloudTextFileUrl)
	secondRetrievalTime := time.Since(start)

	t.Logf("CachedCloudTextFile took %s", secondRetrievalTime)

	assert.NoError(t, err)
	assert.True(t, fromCache)
	assert.Greater(t, len(lines), 0)

	oneThousandthOfFirstRetrieval := firstRetrievalTime / 1000
	assert.Less(t, secondRetrievalTime, oneThousandthOfFirstRetrieval)
}
