package nats

import (
	"time"
)

type NatsSubscribeOpts struct {
	StreamName               string
	SubjectName              string
	ConsumerName             string
	TerminateAfterOneMessage bool
	DupeWindow               time.Duration
}

var defaultNatsSubscribeOpts = &NatsSubscribeOpts{
	StreamName:               "default",
	SubjectName:              "default",
	ConsumerName:             "default",
	TerminateAfterOneMessage: false,
	DupeWindow:               time.Hour * 24 * 90,
}

type NatsPublishOptions struct {
	NatsMessageID        string // Meaningful ID to make use of NATS's deduplication feature
	PersistDeduplication bool
	Subject              string
}

var defaultNatsPublishOptions = &NatsPublishOptions{
	NatsMessageID:        "",
	PersistDeduplication: false,
	Subject:              "",
}

func PublishSubject(subject string) func(*NatsPublishOptions) {
	return func(c *NatsPublishOptions) {
		c.Subject = subject
	}
}

func With10MinuteDuplicationWindowSubscribe() func(*NatsSubscribeOpts) {
	return func(s *NatsSubscribeOpts) {
		s.DupeWindow = time.Minute * 10
	}
}

func WithStreamName(streamName string) func(*NatsSubscribeOpts) {
	return func(s *NatsSubscribeOpts) {
		s.StreamName = streamName
	}
}

func SubscribeSubject(subject string) func(*NatsSubscribeOpts) {
	return func(s *NatsSubscribeOpts) {
		s.SubjectName = subject
	}
}

func StopAfterOneMessage() func(*NatsSubscribeOpts) {
	return func(s *NatsSubscribeOpts) {
		s.TerminateAfterOneMessage = true
	}
}

func WithConsumerName(consumerName string) func(*NatsSubscribeOpts) {
	return func(s *NatsSubscribeOpts) {
		s.ConsumerName = consumerName
	}
}

func StreamNameIsSubjectName() func(*NatsSubscribeOpts) {
	return func(s *NatsSubscribeOpts) {
		s.StreamName = s.SubjectName
	}
}
