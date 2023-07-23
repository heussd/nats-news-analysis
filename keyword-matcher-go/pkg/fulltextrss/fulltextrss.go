package fulltextrss

import (
	"encoding/json"
	"fmt"
	"github.com/heussd/nats-news-keyword-matcher.go/internal/config"
	"io"
	"net/http"
	"os"
	"time"
)

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

func init() {
	for status := 0; status != 200; {
		resp, _ := http.Get(config.FullTextRssServer)
		if resp == nil {
			fmt.Printf("Waiting for %s to come up...\n", config.FullTextRssServer)
			time.Sleep(2 * time.Second)
		} else {
			status = resp.StatusCode
		}
	}
}

func RetrieveFullText(url string) RSSFullTextResponse {
	client := &http.Client{}
	var err error

	req, err := http.NewRequest("GET", config.FullTextRssServer+"/extract.php", nil)
	if err != nil {
		panic(err)
	}

	q := req.URL.Query()
	q.Add("url", url)
	req.URL.RawQuery = q.Encode()

	var response *http.Response
	for {
		var err error
		if response, err = client.Do(req); err != nil {
			panic(err)
		}
		if status := response.StatusCode; status != 200 {
			fmt.Printf("HTTP request failed with status code %d, retrying...\n", status)
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}
	defer response.Body.Close()

	var body []byte
	if body, err = io.ReadAll(response.Body); err != nil {
		panic(err)
	}

	if string(body) == "Invalid URL supplied" {
		fmt.Fprintf(os.Stderr, "Fivefilters service indicated an invalid URL: %s", url)
		return RSSFullTextResponse{}
	}

	var result RSSFullTextResponse
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Fprintf(os.Stderr, "Fivefilters service responded with invalid JSON: %s", string(body))
		return RSSFullTextResponse{}
	}

	return result
}
