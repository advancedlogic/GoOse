package goose

import (
	"github.com/advancedlogic/goquery"
	"gopkg.in/fatih/set.v0"
)

// Article is a collection of properties extracted from the HTML body
type Article struct {
	Title           string
	CleanedText     string
	MetaDescription string
	MetaLang        string
	MetaFavicon     string
	MetaKeywords    string
	CanonicalLink   string
	Domain          string
	TopNode         *goquery.Selection
	TopImage        string
	Tags            *set.Set
	Movies          *set.Set
	FinalURL        string
	LinkHash        string
	RawHTML         string
	Doc             *goquery.Document
	//raw_doc
	PublishDate    string
	AdditionalData map[string]string
	Delta          int64
}

// ToString is a simple method to just show the title
// TODO: add more fields and pretty print
func (article *Article) ToString() string {
	out := article.Title
	return out
}
