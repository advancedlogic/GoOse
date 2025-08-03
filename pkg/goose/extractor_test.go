package goose

import (
	"testing"
)

func Test_title(t *testing.T) {
	titleUnmodified := "   foobar this - is it | bla ¿ "
	title := "foobar this - is it"

	c := NewCrawler(GetDefaultConfiguration())

	a, err := c.Crawl("<!DOCTYPE html><html><head><title>"+
		titleUnmodified+"</title></head></html>", "example.com")
	if err != nil {
		t.Error(err)
	}
	if a.TitleUnmodified != titleUnmodified {
		t.Error("`" + titleUnmodified + "` is extracted as `" + a.TitleUnmodified + "`")
	}

	if a.Title != title {
		t.Error("`" + a.Title + "` should be `" + title + "`")
	}

}
