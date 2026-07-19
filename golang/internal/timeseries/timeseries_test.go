package timeseries

import (
	"testing"

	"github.com/heussd/nats-news-analysis/internal/ngrams"
	"github.com/stretchr/testify/assert"
)

func TestAddTimeSeriesData(t *testing.T) {
	// Test case 1: Add valid time series data
	data := []ngrams.NGram{
		{Words: "example", NGram: 1, Source: "test", Language: "en", Timestamp: "2024-06-01T00:00:00Z"},
		{Words: "test", NGram: 1, Source: "test", Language: "en", Timestamp: "2024-06-01T00:00:00Z"},
	}
	err := AddTimeSeriesData(data)
	if err != nil {
		t.Errorf("AddTimeSeriesData() error = %v, wantErr %v", err, false)
	}

	// Test case 2: Add empty time series data
	data = []ngrams.NGram{}
	err = AddTimeSeriesData(data)
	if err != nil {
		t.Errorf("AddTimeSeriesData() error = %v, wantErr %v", err, false)
	}
}

func TestValidateTimestamp(t *testing.T) {
	assert.Error(t, ValidateTimestamp("2024-06-01 00:00:00"))
	assert.Error(t, ValidateTimestamp("2024-06-01T00:00:00"))
	assert.Error(t, ValidateTimestamp("2024-06-01T00:00:00Z+02:00"))
	assert.Error(t, ValidateTimestamp("2024-06-01"))

	assert.NoError(t, ValidateTimestamp("2024-06-01T00:00:00Z"))
	assert.NoError(t, ValidateTimestamp("2024-06-01 00:00:00+02"))
	assert.NoError(t, ValidateTimestamp("2024-06-01 00:00:00-02"))
	assert.NoError(t, ValidateTimestamp("2026-07-17 15:51:17+00"))
	assert.NoError(t, ValidateTimestamp("2026-07-17 16:00:00+00"))
	assert.NoError(t, ValidateTimestamp("2026-07-17 16:10:05+00"))
	assert.NoError(t, ValidateTimestamp("2026-07-17 16:35:00+00"))
	assert.NoError(t, ValidateTimestamp("2026-07-17 16:52:17+00"))
	assert.NoError(t, ValidateTimestamp("2026-07-17 17:30:00+00"))
	assert.NoError(t, ValidateTimestamp("2026-07-17 17:48:27+00"))
	assert.NoError(t, ValidateTimestamp("2026-07-17 17:53:49+00"))
	assert.NoError(t, ValidateTimestamp("2026-07-17 18:27:28+00"))
}
