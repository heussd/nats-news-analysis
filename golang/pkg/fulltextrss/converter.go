package fulltextrss

import "github.com/heussd/nats-news-analysis/internal/model"

func MakeNews(fulltext RSSFullTextResponse) model.News {
	return model.News{
		Title:    fulltext.Title,
		Excerpt:  fulltext.Excerpt,
		Author:   fulltext.Author,
		Language: fulltext.Language,
		URL:      fulltext.Url,
		Content:  fulltext.Content,
		Date:     fulltext.Date,
	}
}
