package model

import "encoding/json"

type News struct {
	Content       string `json:"content"`     // defuddle JSON
	Title         string `json:"title"`       // defuddle JSON
	Description   string `json:"description"` // defuddle JSON
	Domain        string `json:"domain"`      // defuddle JSON
	Language      string `json:"language"`    // defuddle JSON
	DatePublished string `json:"published"`   // defuddle JSON
	Author        string `json:"author"`      // defuddle JSON
	Site          string `json:"site"`        // defuddle JSON
	WordCount     int    `json:"wordCount"`   // defuddle JSON

	URL string `json:"url"`

	ContentRawText string `json:"contentRawText"`
}

func (m News) GetUrl() string {
	return m.URL
}

// MarshalJSON drops Content: it is only ever read in, never written out.
func (m News) MarshalJSON() ([]byte, error) {
	type alias News
	return json.Marshal(struct {
		alias
		Content string `json:"content,omitempty"`
	}{alias: alias(m)})
}
