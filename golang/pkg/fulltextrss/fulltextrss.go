package fulltextrss

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/heussd/nats-news-analysis/pkg/utils"
)

var (
	FullTextRssServer = utils.GetEnv("FULLTEXTRSS_SERVER", "http://localhost:80")
)

func init() {
	for status := 0; status != 200; {
		resp, _ := http.Get(FullTextRssServer)
		if resp == nil {
			fmt.Printf("Waiting for %s to come up...\n", FullTextRssServer)
			time.Sleep(2 * time.Second)
		} else {
			status = resp.StatusCode
		}
	}
}

func RetrieveFullText(url string) RSSFullTextResponse {
	client := &http.Client{}
	var err error

	req, err := http.NewRequest("GET", FullTextRssServer+"/extract.php", nil)
	if err != nil {
		panic(err)
	}

	q := req.URL.Query()
	q.Add("url", url)
	req.URL.RawQuery = q.Encode()

	var response *http.Response
	for {
		var err error
		status := 0 // If not overwritten indicates failure

		if response, err = client.Do(req); err == nil {
			status = response.StatusCode
			defer response.Body.Close()
			if status == 200 {
				break
			}
		}

		fmt.Printf("HTTP request failed with status code %d, retrying...\n", status)
		time.Sleep(2 * time.Second)
	}

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
