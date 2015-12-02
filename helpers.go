package goose

import (
	"crypto/md5"
	"fmt"
	"github.com/bjarneh/latinx"
	"io"
	"strings"
	"time"
	"unicode/utf8"
)

// Helper is a utility struct to clean up URLs and charsets
type Helper struct {
	urlString string
	url       string
	linkHash  string
}

// NewRawHelper converts the text to UTF8
func NewRawHelper(url string, RawHTML string) Helper {
	if !utf8.ValidString(RawHTML) {
		converter := latinx.Get(latinx.ISO_8859_1)
		RawHTMLBytes, err := converter.Decode([]byte(RawHTML))
		if err != nil {
			fmt.Println(err.Error())
		}
		RawHTML = string(RawHTMLBytes)
	}
	return NewHelper(url)
}

// NewURLHelper wraps the URL
func NewURLHelper(url string) Helper {
	if strings.Contains(url, "#!") {
		return NewHelper(strings.Replace(url, "#!", "?_escaped_fragment_=", -1))
	}
	return NewHelper(url)
}

// NewHelper returns a new Helper
func NewHelper(finalURL string) Helper {
	h := md5.New()
	io.WriteString(h, finalURL)
	bytes := h.Sum(nil)
	helper := Helper{
		urlString: finalURL,
		url:       finalURL,
		linkHash:  fmt.Sprintf("%s.%d", string(bytes), time.Now().UnixNano()),
	}
	return helper
}
