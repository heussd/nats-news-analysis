package likemodel

import "testing"

func TestMatch(t *testing.T) {
	tests := []struct {
		text     string
		expected bool
	}{
		{"some text", false},
		{"podman", true},
		{"lorem lipsum", false},
	}

	for _, tt := range tests {
		result, err := Classify(tt.text)
		if err != nil {
			t.Errorf("Classify(%q) returned error: %v", tt.text, err)
		}
		if result != tt.expected {
			t.Errorf("Classify(%q) = %v; want %v", tt.text, result, tt.expected)
		}
	}
}
