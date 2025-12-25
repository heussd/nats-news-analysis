package nats

import (
	"time"

	"github.com/heussd/nats-news-analysis/internal/model"
)

type NatsStreamOpts struct {
	StreamName  string
	SubjectName string
	DupeWindow  time.Duration
}

var defaultDupeWindow = time.Hour * 24 * 90

type StreamName string

const (
	StreamNameFeedUrls    StreamName = "feed-urls"    // Feeds urls
	StreamNameArticleUrls StreamName = "article-urls" // Articles urls from feeds
	StreamNameNews        StreamName = "news"         // News articles extracted from article urls
	StreamNameMatchUrls   StreamName = "match-urls"   // Urls that are matching interests
)

var StreamConfigs = map[StreamName]NatsStreamOpts{
	StreamNameFeedUrls: {
		StreamName:  string(StreamNameFeedUrls),
		SubjectName: string(StreamNameFeedUrls),
		DupeWindow:  time.Hour * 2,
	},
	StreamNameArticleUrls: {
		StreamName:  string(StreamNameArticleUrls),
		SubjectName: string(StreamNameArticleUrls),
		DupeWindow:  defaultDupeWindow,
	},
	StreamNameNews: {
		StreamName:  string(StreamNameNews),
		SubjectName: "news.*",
		DupeWindow:  defaultDupeWindow,
	},
	StreamNameMatchUrls: {
		StreamName:  string(StreamNameMatchUrls),
		SubjectName: string(StreamNameMatchUrls),
		DupeWindow:  defaultDupeWindow,
	},
}

type NatsSubscribeOpts struct {
	NatsStreamOpts
	ConsumerName             string
	TerminateAfterOneMessage bool
}

type NatsPublishOpts struct {
	NatsMessageID        string // Meaningful ID to make use of NATS's deduplication feature
	PersistDeduplication bool
	Subject              string
}

var defaultNatsPublishOpts = &NatsPublishOpts{
	NatsMessageID:        "",
	PersistDeduplication: false,
	Subject:              "",
}

func PublishSubject(subject string) func(*NatsPublishOpts) {
	return func(c *NatsPublishOpts) {
		c.Subject = subject
	}
}

func StopAfterOneMessage() func(*NatsSubscribeOpts) {
	return func(s *NatsSubscribeOpts) {
		s.TerminateAfterOneMessage = true
	}
}

func SubscribeConsumer(consumerName string) func(*NatsSubscribeOpts) {
	return func(s *NatsSubscribeOpts) {
		s.ConsumerName = consumerName
	}
}

func SubscribeSubject(subjectName string) func(*NatsSubscribeOpts) {
	return func(s *NatsSubscribeOpts) {
		s.SubjectName = subjectName
	}
}

func SubscribeStream(streamName string) func(*NatsSubscribeOpts) {
	return func(s *NatsSubscribeOpts) {
		s.NatsStreamOpts.StreamName = streamName
	}
}

func GetDefaultStreamConfig[T model.PayloadTypes]() *NatsSubscribeOpts {
	var stream NatsStreamOpts
	var v any = new(T)
	switch v.(type) {
	case *model.Article:
		stream = StreamConfigs[StreamNameArticleUrls]
	case *model.Feed:
		stream = StreamConfigs[StreamNameFeedUrls]
	case *model.News:
		stream = StreamConfigs[StreamNameNews]
	case *model.Match:
		stream = StreamConfigs[StreamNameMatchUrls]
	default:
		panic("No stream configuration found for model type")
	}

	return &NatsSubscribeOpts{
		NatsStreamOpts:           stream,
		ConsumerName:             "default",
		TerminateAfterOneMessage: false,
	}
}
