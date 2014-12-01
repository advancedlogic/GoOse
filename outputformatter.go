package goose

import (
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"strconv"
	"strings"
)

type outputFormatter struct {
	topNode  *goquery.Selection
	config   configuration
	language string
}

func (this *outputFormatter) getLanguage(article *Article) string {
	if this.config.useMetaLanguage {
		if article.MetaLang != "" {
			return article.MetaLang
		}
	}
	return this.config.targetLanguage
}

func (this *outputFormatter) getTopNode() *goquery.Selection {
	return this.topNode
}

func (this *outputFormatter) getFormattedText(article *Article) string {
	this.topNode = article.TopNode
	this.language = this.getLanguage(article)
	if this.language == "" {
		this.language = this.config.targetLanguage
	}
	this.removeNegativescoresNodes()
	this.linksToText()
	this.replaceTagsWithText()
	this.removeParagraphsWithFewWords()
	return this.getOutputText()
}

func (this *outputFormatter) convertToText() string {
	txts := make([]string, 0)
	selections := this.topNode
	selections.Each(func(i int, s *goquery.Selection) {
		txt := s.Text()
		if txt != "" {
			txt = txt //unescape
			txtLis := strings.Trim(txt, "\n")
			txts = append(txts, txtLis)
		}
	})
	return strings.Join(txts, "\n\n")
}

func (this *outputFormatter) linksToText() {
	links := this.topNode.Find("a")
	links.Each(func(i int, a *goquery.Selection) {
		imgs := a.Find("img")
		if imgs.Length() == 0 {
			node := a.Get(0)
			node.Data = a.Text()
			node.Type = html.TextNode
		}
	})
}

func (this *outputFormatter) getOutputText() string {
	sb := []string{}
	nodes := this.topNode.Find("*")
	nodes.Each(func(i int, s *goquery.Selection) {
		tag := s.Get(0).DataAtom
		if tag == atom.P {
			sb = append(sb, s.Text())
			sb = append(sb, "\n\n")
		}
	})
	out := strings.Join(sb, "")
	return strings.TrimSpace(out)
}

func (this *outputFormatter) removeNegativescoresNodes() {
	gravityItems := this.topNode.Find("*[gravityScore]")
	gravityItems.Each(func(i int, s *goquery.Selection) {
		score := 0
		sscore, exists := s.Attr("gravityScore")
		if exists {
			score, _ = strconv.Atoi(sscore)
			if score < 1 {
				sNode := s.Get(0)
				sNode.Parent.RemoveChild(sNode)
			}
		}

	})
}

func (this *outputFormatter) replaceTagsWithText() {
	strongs := this.topNode.Find("strong")
	strongs.Each(func(i int, strong *goquery.Selection) {
		text := strong.Text()
		node := strong.Get(0)
		node.Type = html.TextNode
		node.Data = text
	})

	bolds := this.topNode.Find("b")
	bolds.Each(func(i int, bold *goquery.Selection) {
		text := bold.Text()
		node := bold.Get(0)
		node.Type = html.TextNode
		node.Data = text
	})

	italics := this.topNode.Find("i")
	italics.Each(func(i int, italic *goquery.Selection) {
		text := italic.Text()
		node := italic.Get(0)
		node.Type = html.TextNode
		node.Data = text
	})
}

func (this *outputFormatter) removeParagraphsWithFewWords() {
	language := this.language
	if language == "" {
		language = "en"
	}
	allNodes := this.topNode.Find("*")
	allNodes.Each(func(i int, s *goquery.Selection) {
		sw := this.config.stopWords.stopWordsCount(language, s.Text())
		if sw.stopWordCount < 3 && s.Find("object").Length() == 0 &&
			s.Find("object").Length() == 0 && s.Find("em").Length() == 0 {
			node := s.Get(0)
			node.Parent.RemoveChild(node)
		}
	})

}
