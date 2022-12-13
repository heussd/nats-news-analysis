package fulltextrss

import (
	"encoding/json"
	"github.com/heussd/nats-news-keyword-matcher.go/internal/utils"
	"io"
	"net/http"
	"strconv"
)

var server = utils.GetEnv("FULLTEXTRSS_SERVER", "http://localhost:80")

type RSSFullTextResponse struct {
	Title              string `json:"title"`
	Excerpt            string `json:"excerpt"`
	Date               string `json:"date"`
	Author             string `json:"author"`
	Language           string `json:"language"`
	Url                string `json:"url""`
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

func RetrieveFullText(url string) RSSFullTextResponse {
	client := &http.Client{}
	req, err := http.NewRequest("GET", server+"/extract.php", nil)
	if err != nil {
		panic(err)
	}

	q := req.URL.Query()
	q.Add("url", url)
	req.URL.RawQuery = q.Encode()

	response, error := client.Do(req)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	if status := response.StatusCode; status != 200 {
		panic("HTTP request failed with status code " + strconv.Itoa(status))
	}

	body, _ := io.ReadAll(response.Body)

	var result RSSFullTextResponse
	if err := json.Unmarshal(body, &result); err != nil {
		panic(err)
	}

	return result
}
