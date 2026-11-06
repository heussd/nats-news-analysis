package defuddle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/heussd/nats-news-analysis/internal/htmlsanitise"
	"github.com/heussd/nats-news-analysis/internal/model"
	"github.com/heussd/nats-news-analysis/pkg/utils"
)

var DefuddleServer = utils.GetEnv("DEFUDDLE_SERVER", "http://loadbalancer:80")

func init() {
	fmt.Printf("Using %s...\n", DefuddleServer)
	for status := 0; status != 200; {
		resp, _ := http.Get(DefuddleServer)
		if resp == nil {
			fmt.Printf("Waiting for %s to come up...\n", DefuddleServer)
			time.Sleep(2 * time.Second)
		} else {
			status = resp.StatusCode
		}
	}
}

func RetrieveFullText(url string) (news model.News, err error) {
	client := &http.Client{}

	jsonData := map[string]interface{}{
		"url": url,
		"options": map[string]bool{
			"removeHiddenElements": true,
			"removeLowScoring":     true,
			"removeSmallImages":    true,
			"removeImages":         true,
			"standardize":          true,
			"useAsync":             true,
		},
	}
	jsonBytes, _ := json.Marshal(jsonData)

	req, err := http.NewRequest("POST", DefuddleServer+"/parse", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}

	var response *http.Response
	status := 0 // If not overwritten indicates failure
	backoff := 2 * time.Second
	const maxBackoff = 10 * time.Second
	for {
		var err error

		if response, err = client.Do(req); err == nil {
			status = response.StatusCode
			defer response.Body.Close()
			if status == 200 {
				break
			}
		}

		fmt.Printf("HTTP request failed with status code %d, retrying in %s...\n", status, backoff)
		time.Sleep(backoff)
		if backoff *= 2; backoff > maxBackoff {
			break
		}
	}

	if status != 200 {
		return news, fmt.Errorf("failed to retrieve full text, status code: %d", status)
	}

	var body []byte
	if body, err = io.ReadAll(response.Body); err != nil {
		return news, fmt.Errorf("failed to read response body: %w", err)
	}

	if err = json.Unmarshal(body, &news); err != nil {
		return news, fmt.Errorf("invalid JSON response")
	}

	news.ContentRawText = htmlsanitise.Sanitize(news.Content)

	news.URL = url
	return news, nil
}
