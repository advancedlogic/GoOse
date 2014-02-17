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
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

type Crawler struct {
	config  configuration
	url     string
	rawHtml string
	helper  Helper
}

func NewCrawler(config configuration, url string, rawHtml string) Crawler {
	return Crawler{
		config:  config,
		url:     url,
		rawHtml: rawHtml,
	}
}

func (this Crawler) Crawl() *Article {
	article := new(Article)
	this.assignParseCandidate()
	this.assignHtml()

	if this.rawHtml == "" {
		return article
	}

	reader := strings.NewReader(this.rawHtml)
	document, err := goquery.NewDocumentFromReader(reader)

	if err == nil {
		extractor := NewExtractor(this.config)
		html, _ := document.Html()
		start := TimeInNanoseconds()
		article.RawHtml = html
		article.FinalUrl = this.helper.url
		article.LinkHash = this.helper.linkHash
		article.Doc = document
		article.Title = extractor.getTitle(article)
		article.MetaLang = extractor.getMetaLanguage(article)
		article.MetaFavicon = extractor.getFavicon(article)

		article.MetaDescription = extractor.getMetaContentWithSelector(article, "meta[name=description]")
		article.MetaKeywords = extractor.getMetaContentWithSelector(article, "meta[name=keywords]")
		article.CanonicalLink = extractor.getCanonicalLink(article)
		article.Domain = extractor.getDomain(article)
		article.Tags = extractor.getTags(article)

		cleaner := NewCleaner(this.config)
		article.Doc = cleaner.clean(article)

		article.TopNode = extractor.calculateBestNode(article)
		if article.TopNode != nil {
			article.TopNode = extractor.postCleanup(article.TopNode)

			outputFormatter := new(outputFormatter)
			article.CleanedText = outputFormatter.getFormattedText(article)

			videoExtractor := NewVideoExtractor()
			article.Movies = videoExtractor.GetVideos(article)

			article.TopImage = OpenGraphResolver(article)
			if article.TopImage == "" {
				article.TopImage = WebPageResolver(article)
			}
		}

		stop := TimeInNanoseconds()
		delta := stop - start
		article.Delta = delta

	} else {
		panic(err.Error())
	}
	return article
}

func (this *Crawler) assignParseCandidate() {
	if this.rawHtml != "" {
		this.helper = NewRawHelper(this.url, this.rawHtml)
	} else {
		this.helper = NewUrlHelper(this.url)
	}
}

func (this *Crawler) assignHtml() {
	if this.rawHtml == "" {
		cookieJar, _ := cookiejar.New(nil)
		client := &http.Client{
			Jar: cookieJar,
		}
		req, err := http.NewRequest("GET", this.url, nil)
		if err == nil {
			req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_6_7) AppleWebKit/534.30 (KHTML, like Gecko) Chrome/12.0.742.91 Safari/534.30")
			resp, err := client.Do(req)
			if err == nil {
				defer resp.Body.Close()
				contents, err := ioutil.ReadAll(resp.Body)
				if err == nil {
					this.rawHtml = string(contents)
				} else {
					log.Println(err.Error())
				}
			} else {
				log.Println(err.Error())
			}
		} else {
			log.Println(err.Error())
		}
	}
}
