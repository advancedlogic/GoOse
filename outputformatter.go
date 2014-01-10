package goose

import (
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

type OutputFormatter struct {
	topNode  *goquery.Selection
	config   Configuration
	language string
	parser   Parser
}

func (of *OutputFormatter) getLanguage(article *Article) string {
	if of.config.UseMetaLanguage {
		if article.MetaLang != "" {
			return article.MetaLang
		}
	}
	return of.config.TargetLanguage
}

func (of *OutputFormatter) getTopNode() *goquery.Selection {
	return of.topNode
}

func (of *OutputFormatter) getFormattedText(article *Article) string {
	of.topNode = article.TopNode
	of.language = of.getLanguage(article)
	of.removeNegativescoresNodes()
	of.linksToText()
	of.replaceTagsWithText()
	of.removeParagraphsWithFewWords()

	return of.getOutputText()
}

func (of *OutputFormatter) convertToText() string {
	txts := make([]string, 0)
	selections := of.topNode
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

func (of *OutputFormatter) linksToText() {
	links := of.topNode.Find("a")
	links.Each(func(i int, a *goquery.Selection) {
		imgs := a.Find("img")
		if imgs.Length() == 0 {
			node := a.Get(0)
			text := a.Text()
			node.Type = html.TextNode
			node.Data = text
		}
	})
}

func (of *OutputFormatter) getOutputText() string {
	sb := make([]string, 0)
	nodes := of.topNode.Find("*")
	nodes.Each(func(i int, s *goquery.Selection) {
		tag := s.Get(0).DataAtom
		if tag == atom.P {
			text := s.Text()
			sb = append(sb, text)
			sb = append(sb, "\n\n")
		}
	})
	return strings.Join(sb, "")
}

func (of *OutputFormatter) removeNegativescoresNodes() {
	gravityItems := of.topNode.Find("*[gravityScore]")
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

func (of *OutputFormatter) replaceTagsWithText() {
	strongs := of.topNode.Find("strong")
	strongs.Each(func(i int, strong *goquery.Selection) {
		text := strong.Text()
		node := strong.Get(0)
		node.Type = html.TextNode
		node.Data = text
	})

	bolds := of.topNode.Find("b")
	bolds.Each(func(i int, bold *goquery.Selection) {
		text := bold.Text()
		node := bold.Get(0)
		node.Type = html.TextNode
		node.Data = text
	})

	italics := of.topNode.Find("i")
	italics.Each(func(i int, italic *goquery.Selection) {
		text := italic.Text()
		node := italic.Get(0)
		node.Type = html.TextNode
		node.Data = text
	})
}

func (of *OutputFormatter) removeParagraphsWithFewWords() {
	language := of.language
	if language == "" {
		language = "en"
	}
	allNodes := of.topNode.Find("*")
	allNodes.Each(func(i int, s *goquery.Selection) {
		stopWordsCount := gooseStopWordsCount(language, s.Text())
		if stopWordsCount < 5 && s.Find("object").Length() == 0 && s.Find("embed").Length() == 0 {
			node := s.Get(0)
			node.Parent.RemoveChild(node)
		}
	})

}
