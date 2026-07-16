package timeseries

import (
	"testing"

	"github.com/heussd/nats-news-analysis/internal/ngrams"
)

func TestAddTimeSeriesData(t *testing.T) {
	// Test case 1: Add valid time series data
	data := []ngrams.NGram{
		{Words: "example", Count: 5, NGram: 1, Source: "test", Language: "en", Timestamp: "2024-06-01T00:00:00Z"},
		{Words: "test", Count: 3, NGram: 1, Source: "test", Language: "en", Timestamp: "2024-06-01T00:00:00Z"},
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
