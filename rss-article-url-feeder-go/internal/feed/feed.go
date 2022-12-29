package feed

import (
	"fmt"
	"github.com/SlyMarbo/rss"
	"net/url"
	"strings"
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

func ArticleUrls(feedUrl string) []string {
	var articleUrls []string

	feed, err := rss.Fetch(feedUrl)
	if err != nil {
		fmt.Printf("❌️ %s failed to load: %s\n", feedUrl, err)
		return articleUrls
	}

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
	return articleUrls
}
