/*
This is a golang port of "Goose" originaly licensed to Gravity.com
under one or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.

Golang port was written by Antonio Linari

Gravity.com licenses this file
to you under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package goose

import (
	"crypto/md5"
	"fmt"
	"github.com/bjarneh/latinx"
	"io"
	"strings"
	"unicode/utf8"
)

type Helper struct {
	urlString string
	url       string
	linkHash  string
}

func NewRawHelper(url string, rawHtml string) Helper {
	if utf8.ValidString(rawHtml) {
		converter := latinx.Get(latinx.ISO_8859_1)
		rawHtmlBytes, err := converter.Decode([]byte(rawHtml))
		if err != nil {
			fmt.Println(err.Error())
		}
		rawHtml = string(rawHtmlBytes)
	}
	h := md5.New()
	io.WriteString(h, url)
	bytes := h.Sum(nil)
	helper := Helper{
		urlString: url,
		url:       url,
		linkHash:  fmt.Sprintf("%s.%d", string(bytes), TimeInNanoseconds()),
	}
	return helper
}

func NewUrlHelper(url string) Helper {
	finalUrl := ""
	if strings.Contains(url, "#!") {
		finalUrl = strings.Replace(url, "#!", "?_escaped_fragment_=", -1)
	} else {
		finalUrl = url
	}
	h := md5.New()
	io.WriteString(h, finalUrl)
	bytes := h.Sum(nil)
	helper := Helper{
		urlString: finalUrl,
		url:       finalUrl,
		linkHash:  fmt.Sprintf("%s.%d", string(bytes), TimeInNanoseconds()),
	}
	return helper
}
