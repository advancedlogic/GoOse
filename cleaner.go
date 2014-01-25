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
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"container/list"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
)

type cleaner struct {
}

func NewCleaner() cleaner {
	return cleaner{}
}

var divToPElementsPattern = regexp.MustCompile("<(a|blockquote|dl|div|img|ol|p|pre|table|ul)")
var tabsRegEx, _ = regexp.Compile("\\t|^\\s+$]")
var REMOVENODES_RE = regexp.MustCompile("^side$|combx|retweet|mediaarticlerelated|menucontainer|navbar|comment|PopularQuestions|contact|foot|footer|Footer|footnote|cnn_strycaptiontxt|cnn_html_slideshow|cnn_strylftcntnt|links|meta$|scroll|shoutbox|sponsor|tags|socialnetworking|socialNetworking|cnnStryHghLght|cnn_stryspcvbx|^inset$|pagetools|post-attributes|welcome_form|contentTools2|the_answers|communitypromo|runaroundLeft|subscribe|vcard|articleheadings|date|^print$|popup|author-dropdown|tools|socialtools|byline|konafilter|KonaFilter|breadcrumbs|^fn$|wp-caption-text|legende|ajoutVideo|timestamp|js_replies")
var CAPTIONS_RE = regexp.MustCompile("^caption$")
var GOOGLE_RE = regexp.MustCompile(" google ")
var MORE_RE = regexp.MustCompile("^[^entry-]more.*$")
var FACEBOOK_RE = regexp.MustCompile("[^-]facebook")
var FACEBOOK_BROADCASTING_RE = regexp.MustCompile("facebook-broadcasting")
var TWITTER_RE = regexp.MustCompile("[^-]twitter")

func (this *cleaner) clean(article *Article) *goquery.Document {
	docToClean := article.Doc
	docToClean = this.cleanArticleTags(docToClean)
	docToClean = this.cleanEMTags(docToClean)
	docToClean = this.dropCaps(docToClean)
	docToClean = this.removeScriptsStyle(docToClean)
	docToClean = this.removeNodesRegEx(docToClean, REMOVENODES_RE)
	docToClean = this.removeNodesRegEx(docToClean, CAPTIONS_RE)
	docToClean = this.removeNodesRegEx(docToClean, GOOGLE_RE)
	docToClean = this.removeNodesRegEx(docToClean, MORE_RE)
	docToClean = this.removeNodesRegEx(docToClean, FACEBOOK_RE)
	docToClean = this.removeNodesRegEx(docToClean, FACEBOOK_BROADCASTING_RE)
	docToClean = this.removeNodesRegEx(docToClean, TWITTER_RE)

	docToClean = this.cleanParaSpans(docToClean)
	docToClean = this.convertDivsToParagraphs(docToClean, "div")
	docToClean = this.convertDivsToParagraphs(docToClean, "span")
	docToClean = this.convertDivsToParagraphs(docToClean, "article")
	return docToClean
}

func (this *cleaner) cleanArticleTags(doc *goquery.Document) *goquery.Document {
	tags := [3]string{"id", "name", "class"}
	articles := doc.Find("article")
	articles.Each(func(i int, s *goquery.Selection) {
		for _, tag := range tags {
			this.delAttr(s, tag)
		}
	})
	return doc
}

func (this *cleaner) cleanEMTags(doc *goquery.Document) *goquery.Document {
	ems := doc.Find("em")
	ems.Each(func(i int, s *goquery.Selection) {
		images := s.Find("img")
		if images.Length() == 0 {
			node := s.Get(0)
			node.Data = s.Text()
			node.Type = html.TextNode
		}
	})
	return doc
}

func (this *cleaner) cleanCites(doc *goquery.Document) *goquery.Document {
	cites := doc.Find("cite")
	cites.Each(func(i int, s *goquery.Selection) {
		node := s.Get(0)
		node.Parent.RemoveChild(node)
	})
	return doc
}

func (this *cleaner) cleanDivs(doc *goquery.Document) *goquery.Document {
	frames := make(map[string]int)
	framesNodes := make(map[string]*list.List)
	divs := doc.Find("div")
	divs.Each(func(i int, s *goquery.Selection) {
		children := s.Children()
		if children.Size() == 0 {
			text := s.Text()
			text = strings.Trim(text, " ")
			text = strings.Trim(text, "\t")
			text = strings.ToLower(text)
			frames[text]++
			if framesNodes[text] == nil {
				framesNodes[text] = list.New()
			}
			framesNodes[text].PushBack(s)
		}
	})
	for text, freq := range frames {
		if freq > 1 {
			selections := framesNodes[text]
			for s := selections.Front(); s != nil; s = s.Next() {
				selection := s.Value.(*goquery.Selection)
				node := selection.Get(0)
				node.Parent.RemoveChild(node)
			}
		}
	}
	return doc
}

func (this *cleaner) dropCaps(doc *goquery.Document) *goquery.Document {
	items := doc.Find("span")
	items.Each(func(i int, s *goquery.Selection) {
		attribute, exists := s.Attr("class")
		if exists && (strings.Contains(attribute, "dropcap") || strings.Contains(attribute, "drop_cap")) {
			node := s.Get(0)
			node.Type = html.TextNode
			node.Data = s.Text()
		}
	})
	return doc
}

func (this *cleaner) removeScriptsStyle(doc *goquery.Document) *goquery.Document {
	scripts := doc.Find("script")
	scripts.Each(func(i int, s *goquery.Selection) {
		node := s.Get(0)
		node.Parent.RemoveChild(node)
	})

	styles := doc.Find("style")
	styles.Each(func(i int, s *goquery.Selection) {
		node := s.Get(0)
		node.Parent.RemoveChild(node)
	})

	//remove comments :) How????
	return doc
}

func (this *cleaner) matchNodeRegEx(attribute string, pattern *regexp.Regexp) bool {
	return pattern.MatchString(attribute)
}

func (this *cleaner) removeNodesRegEx(doc *goquery.Document, pattern *regexp.Regexp) *goquery.Document {
	//println("removeNodesRegEx")
	selectors := [3]string{"id", "class", "name"}
	for _, selector := range selectors {
		naughtyList := doc.Find("*")
		naughtyList.Each(func(i int, s *goquery.Selection) {
			attribute, _ := s.Attr(selector)
			if this.matchNodeRegEx(attribute, pattern) {
				node := s.Get(0)
				node.Parent.RemoveChild(node)
			}
		})
	}
	return doc
}

func (this *cleaner) cleanParaSpans(doc *goquery.Document) *goquery.Document {
	spans := doc.Find("span")
	spans.Each(func(i int, s *goquery.Selection) {
		if s.Parent().Get(0).DataAtom == atom.P {
			node := s.Get(0)
			node.Data = s.Text()
			node.Type = html.TextNode
		}
	})
	return doc
}

func (this *cleaner) getFlushedBuffer(replacementText string) *goquery.Selection {
	reader := strings.NewReader(replacementText)
	document, err := goquery.NewDocumentFromReader(reader)
	if err == nil {
		return document.Selection
	}
	return nil
}

func (this *cleaner) replaceWithPara(div *goquery.Selection) {
	node := div.Get(0)
	node.Data = div.Text()
	node.Type = html.ElementNode
	node.DataAtom = atom.P
}

func (this *cleaner) convertDivsToParagraphs(doc *goquery.Document, domType string) *goquery.Document {
	badDivs := 0
	convertedTextNodes := 0
	divs := doc.Find(domType)

	divs.Each(func(i int, div *goquery.Selection) {
		h, _ := div.Html()
		h = strings.ToLower(h)

		if divToPElementsPattern.Match([]byte(h)) {
			this.replaceWithPara(div)
			badDivs++
		} else {
			replacementText := make([]string, 0)
			nodesToRemove := list.New()

			div.Siblings().EachWithBreak(func(i int, kid *goquery.Selection) bool {
				kidNode := kid.Get(0)
				if kidNode != nil {
					child := kidNode.FirstChild
					for child != nil {
						text := child.Data
						if text == "" {
							return false
						}
						if child.DataAtom == atom.P && len(replacementText) > 0 {
							node := new(html.Node)
							node.DataAtom = atom.P
							node.Data = kid.Text()
							node.Type = html.ElementNode
							div.Children().AddNodes(node)
							replacementText = make([]string, 0)
						} else if child.Type == html.TextNode {
							text = strings.Replace(text, "\n", "", -1)
							text = tabsRegEx.ReplaceAllString(text, "")
							if len(text) > 1 {
								prev := child.PrevSibling
								if prev != nil && prev.DataAtom == atom.A {
									nodeSelection := kid.HasNodes(prev)
									html, _ := nodeSelection.Siblings().Html()
									replacementText = append(replacementText, html)
								}
								replacementText = append(replacementText, text)
								nodesToRemove.PushBack(child)
								convertedTextNodes++
							}

						}
						child = child.NextSibling
					}
				}
				return true
			})

			newNode := new(html.Node)
			newNode.Type = html.ElementNode
			newNode.Data = strings.Join(replacementText, "")
			newNode.DataAtom = atom.P
			div.First().AddNodes(newNode)

			for s := nodesToRemove.Front(); s != nil; s = s.Next() {
				node := s.Value.(*html.Node)
				if node != nil && node.Parent != nil {
					node.Parent.RemoveChild(node)
				}
			}
		}
	})

	return doc

}

func (this *cleaner) indexOfAttribute(selection *goquery.Selection, attr string) int {
	node := selection.Get(0)
	for i, a := range node.Attr {
		if a.Key == attr {
			return i
		}
	}
	return -1
}

func (this *cleaner) delAttr(selection *goquery.Selection, attr string) {
	idx := this.indexOfAttribute(selection, attr)
	if idx > -1 {
		node := selection.Get(0)
		node.Attr = append(node.Attr[:idx], node.Attr[idx+1:]...)
	}
}
