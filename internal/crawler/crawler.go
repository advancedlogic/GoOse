package crawler

import (
	"errors"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/advancedlogic/GoOse/internal/extractor"
	"github.com/advancedlogic/GoOse/internal/utils"
	"github.com/advancedlogic/GoOse/pkg/goose"
)

// Crawler can fetch the target HTML page
type Crawler struct {
	config  goose.Configuration
	Charset string
}

// NewCrawler returns a crawler object initialised with the URL and the [optional] raw HTML body
func NewCrawler(config goose.Configuration) Crawler {
	return Crawler{
		config:  config,
		Charset: "",
	}
}

func getCharsetFromContentType(cs string) string {
	cs = strings.ToLower(strings.Replace(cs, " ", "", -1))
	cs = strings.TrimPrefix(cs, "text/html;charset=")
	cs = strings.TrimPrefix(cs, "text/xhtml;charset=")
	cs = strings.TrimPrefix(cs, "application/xhtml+xml;charset=")
	return utils.NormaliseCharset(cs)
}

// SetCharset can be used to force a charset (e.g. when read from the HTTP headers)
// rather than relying on the detection from the HTML meta tags
func (c *Crawler) SetCharset(cs string) {
	c.Charset = getCharsetFromContentType(cs)
}

// GetContentType returns the Content-Type string extracted from the meta tags
func (c Crawler) GetContentType(document *goquery.Document) string {
	var attr string
	// <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	document.Find("meta[http-equiv#=(?i)^Content\\-type$]").Each(func(i int, s *goquery.Selection) {
		attr, _ = s.Attr("content")
	})
	return attr
}

// GetCharset returns a normalised charset string extracted from the meta tags
func (c Crawler) GetCharset(document *goquery.Document) string {
	// manually-provided charset (from HTTP headers?) takes priority
	if "" != c.Charset {
		return c.Charset
	}

	// <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	ct := c.GetContentType(document)
	if "" != ct && strings.Contains(strings.ToLower(ct), "charset") {
		return getCharsetFromContentType(ct)
	}

	// <meta charset="utf-8">
	selection := document.Find("meta").EachWithBreak(func(i int, s *goquery.Selection) bool {
		_, exists := s.Attr("charset")
		return !exists
	})

	if selection != nil {
		cs, _ := selection.Attr("charset")
		return utils.NormaliseCharset(cs)
	}

	return ""
}

// Preprocess fetches the HTML page if needed, converts it to UTF-8 and applies
// some text normalisation to guarantee better results when extracting the content
func (c *Crawler) Preprocess(RawHTML string) (*goquery.Document, error) {
	var err error

	if RawHTML == "" {
		return nil, errors.New("cannot process empty HTML content")
	}

	RawHTML = c.addSpacesBetweenTags(RawHTML)

	reader := strings.NewReader(RawHTML)
	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	cs := c.GetCharset(document)
	//log.Println("-------------------------------------------CHARSET:", cs)
	if "" != cs && "UTF-8" != cs {
		// the net/html parser and goquery require UTF-8 data
		RawHTML = utils.UTF8encode(RawHTML, cs)
		reader = strings.NewReader(RawHTML)
		if document, err = goquery.NewDocumentFromReader(reader); err != nil {
			return nil, err
		}
	}

	return document, nil
}

// Crawl fetches the HTML body and returns an Article
func (c Crawler) Crawl(RawHTML string, url string) (*goose.Article, error) {
	article := new(goose.Article)

	document, err := c.Preprocess(RawHTML)
	if nil != err {
		return nil, err
	}
	if nil == document {
		return article, nil
	}
	extr := extractor.NewExtractor(c.config)
	startTime := time.Now().UnixNano()

	article.RawHTML, err = document.Html()
	if nil != err {
		return nil, err
	}
	article.FinalURL = url
	article.Doc = document

	article.Title = extr.GetTitle(document)
	article.TitleUnmodified = article.Title
	article.MetaLang = extr.GetMetaLanguage(document)
	article.MetaFavicon = extr.GetFavicon(document)

	article.MetaDescription = extr.GetMetaContentWithSelector(document, "meta[name#=(?i)^description$]")
	article.MetaKeywords = extr.GetMetaContentWithSelector(document, "meta[name#=(?i)^keywords$]")
	article.CanonicalLink = extr.GetCanonicalLink(document)
	if "" == article.CanonicalLink {
		article.CanonicalLink = article.FinalURL
	}
	article.Domain = extr.GetDomain(article.CanonicalLink)
	article.Tags = extr.GetTags(document)

	if c.config.ExtractPublishDate {
		if timestamp := extr.GetPublishDate(document); timestamp != nil {
			article.PublishDate = timestamp
		}
	}

	cleaner := extractor.NewCleaner(c.config)
	article.Doc = cleaner.Clean(article.Doc)

	article.TopImage = extractor.OpenGraphResolver(document)
	if article.TopImage == "" {
		article.TopImage = extractor.WebPageResolver(article)
	}

	article.TopNode = extr.CalculateBestNode(document)
	if article.TopNode != nil {
		article.TopNode = extr.PostCleanup(article.TopNode)

		article.CleanedText, article.Links = extr.GetCleanTextAndLinks(article.TopNode, article.MetaLang)

		videoExtractor := extractor.NewVideoExtractor()
		article.Movies = videoExtractor.GetVideos(document)
	}

	article.Delta = time.Now().UnixNano() - startTime

	return article, nil
}

// In many cases, like at the end of each <li> element or between </span><span> tags,
// we need to add spaces, otherwise the text on either side will get joined together into one word.
// This method also adds newlines after each </p> tag to preserve paragraphs.
func (c Crawler) addSpacesBetweenTags(text string) string {
	text = strings.Replace(text, "><", "> <", -1)
	text = strings.Replace(text, "</blockquote>", "</blockquote>\n", -1)
	text = strings.Replace(text, "<img ", "\n<img ", -1)
	text = strings.Replace(text, "</li>", "</li>\n", -1)
	return strings.Replace(text, "</p>", "</p>\n", -1)
}
