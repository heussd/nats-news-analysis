package model

import "github.com/heussd/nats-news-analysis/pkg/fulltextrss"

func MakeNews(fulltext fulltextrss.RSSFullTextResponse) News {
	return News{
		Title:    fulltext.Title,
		Excerpt:  fulltext.Excerpt,
		Author:   fulltext.Author,
		Language: fulltext.Language,
		URL:      fulltext.Url,
		Content:  fulltext.Content,
		Date:     fulltext.Date,
	}
}
