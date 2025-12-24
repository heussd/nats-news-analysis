package cloudtextfile

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	lastGenerated time.Time
	cacheDuration = 20 * time.Minute
	cachedLines   []string
)

func retrieveCloudTexFile(url string) (lines []string, err error) {
	client := &http.Client{}
	var req *http.Request

	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return nil, err
	}

	var response *http.Response

	if response, err = client.Do(req); err != nil {
		return nil, err
	}

	if status := response.StatusCode; status != 200 {
		return nil, fmt.Errorf("failed to retrieve keywords file: status code %d", status)
	}

	var body []byte
	if body, err = io.ReadAll(response.Body); err != nil {
		return nil, err
	}

	lines = strings.Split(string(body), "\n")

	return lines, nil
}

func CachedCloudTextFile(url string) (lines []string, fromCache bool, err error) {
	if time.Since(lastGenerated) < cacheDuration {
		return cachedLines, true, nil
	}

	if cachedLines, err = retrieveCloudTexFile(url); err != nil {
		return nil, false, fmt.Errorf("could not retrieve cloud text file %w", err)
	}

	lastGenerated = time.Now()
	return cachedLines, false, nil
}
