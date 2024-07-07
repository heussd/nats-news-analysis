package feed

import (
	"fmt"
	"github.com/SlyMarbo/rss"
	"net/http"
	"net/url"
	"strings"
	"time"
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

func ArticleUrls(feedUrl string) []string {
	var articleUrls []string

	feed, err := rss.FetchByFunc(TimeoutFetchFunc, feedUrl)
	if err != nil {
		fmt.Printf("❌️ %s failed to load: %s\n", feedUrl, err)
		return articleUrls
	}

	// Force feed update
	feed.Refresh = time.Now()

	err = feed.Update()
	if err != nil {
		fmt.Printf("⚠️ %s not updated: %s\n", feedUrl, err)
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
		fmt.Printf("❌ No articles found in %s\n", feedUrl)
	}

	fmt.Printf("✅ %3d items found in %s\n", len(articleUrls), feedUrl)
	return articleUrls
}
