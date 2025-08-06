package fulltextrss

type RSSFullTextResponse struct {
	Title              string `json:"title"`
	Excerpt            string `json:"excerpt"`
	Date               string `json:"date"`
	Author             string `json:"author"`
	Language           string `json:"language"`
	Url                string `json:"url"`
	EffectiveUrl       string `json:"effective_url"`
	Domain             string `json:"domain"`
	WordCount          int    `json:"word_count"`
	OgUrl              string `json:"og_url"`
	OgTitle            string `json:"og_title"`
	OgDescription      string `json:"og_description"`
	OgImage            string `json:"og_image"`
	OgType             string `json:"og_type"`
	TwitterCard        string `json:"twitter_card"`
	TwitterSite        string `json:"twitter_site"`
	TwitterCreator     string `json:"twitter_creator"`
	TwitterImage       string `json:"twitter_image"`
	TwitterTitle       string `json:"twitter_title"`
	TwitterDescription string `json:"twitter_description"`
	Content            string `json:"content"`
}
