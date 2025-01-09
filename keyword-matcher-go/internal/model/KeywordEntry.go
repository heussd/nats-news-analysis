package model

import "github.com/dlclark/regexp2"

type KeywordEntry struct {
	Regexp regexp2.Regexp
	Id     string
	Text   string
}
