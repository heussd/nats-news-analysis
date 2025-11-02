package model

import "github.com/dlclark/regexp2"

type PayloadTypes interface {
	Article | Match | Feed | News

	GetUrl() string
}

type Article struct {
	Url string `json:"url"`
}

func (a Article) GetUrl() string {
	return a.Url
}

type Feed struct {
	Url string `json:"url"`
}

func (f Feed) GetUrl() string {
	return f.Url
}

type Match struct {
	Keywords []Keyword `json:"keywords"`
	Url      string    `json:"url"`
}

func (m Match) GetUrl() string {
	return m.Url
}

type Keyword struct {
	Regexp regexp2.Regexp `json:"-"` // Don't serialize
	Id     string         `json:"id"`
	Text   string         `json:"Text"` // The original regex text
}

type News struct {
	Title    string `json:"title"`
	Excerpt  string `json:"excerpt"`
	Author   string `json:"author"`
	Language string `json:"language"`
	URL      string `json:"url"`
	Content  string `json:"content"`
	Date     string `json:"date"`
}

func (m News) GetUrl() string {
	return m.URL
}
