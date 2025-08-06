package feed

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/SlyMarbo/rss"
)

// https://stackoverflow.com//questions/72387330/how-to-extract-base-url-using-golang#answer-72387843
func getBaseUrl(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	u.Path = ""
	u.RawQuery = ""
	u.Fragment = ""

	return u.String()
}

var TimeoutFetchFunc = func(url string) (resp *http.Response, err error) {
	var client = &http.Client{
		Timeout: time.Second * 60,
	}
	return client.Get(url)
}

func FetchFeedAndExtractArticleUrls(feedUrl string) (articleUrls []string, err error) {
	var feed *rss.Feed
	if feed, err = rss.FetchByFunc(TimeoutFetchFunc, feedUrl); err != nil {
		return articleUrls, fmt.Errorf("failed to fetch %s: %w", feedUrl, err)
	}

	// Force feed update
	feed.Refresh = time.Now()

	if err = feed.Update(); err != nil {
		return articleUrls, fmt.Errorf("feed %s not updated: %w", feedUrl, err)
	}

	for _, item := range feed.Items {
		link := item.Link

		// Some feeds only serve relative article urls
		if !strings.HasPrefix(link, "http") {
			link = getBaseUrl(feedUrl) + "/" + link
		}

		articleUrls = append(articleUrls, link)
	}

	if len(articleUrls) == 0 {
		return articleUrls, fmt.Errorf("no articles found in feed %s", feedUrl)
	}

	return articleUrls, nil
}
