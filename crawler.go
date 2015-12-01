package goose

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Crawler can fetch the target HTML page
type Crawler struct {
	config  Configuration
	url     string
	RawHTML string
	Charset string
	helper  Helper
}

// NewCrawler returns a crawler object initialised with the URL and the [optional] raw HTML body
func NewCrawler(config Configuration, url string, RawHTML string) Crawler {
	return Crawler{
		config:  config,
		url:     url,
		RawHTML: RawHTML,
		Charset: "",
	}
}

// SetCharset can be used to force a charset (e.g. when read from the HTTP headers)
// rather than relying on the detection from the HTML meta tags
func (c *Crawler) SetCharset(cs string) {
	c.Charset = NormaliseCharset(cs)
}

// GetCharset returns a normalised charset string extracted from the meta tags
func (c Crawler) GetCharset(document *goquery.Document) string {
	// manually-provided charset (from HTTP headers?) takes priority
	if "" != c.Charset {
		return c.Charset
	}

	// <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	selection := document.Find("meta").EachWithBreak(func(i int, s *goquery.Selection) bool {
		attr, exists := s.Attr("http-equiv")
		if exists && attr == "Content-Type" {
			return false
		}
		return true
	})

	if selection != nil {
		attr, _ := selection.Attr("content")
		attr = strings.Replace(attr, " ", "", -1)

		if strings.HasPrefix(attr, "text/html;charset=") {
			cs := strings.TrimPrefix(attr, "text/html;charset=")
			return NormaliseCharset(cs)
		}
	}

	// <meta charset="utf-8">
	selection = document.Find("meta").EachWithBreak(func(i int, s *goquery.Selection) bool {
		_, exists := s.Attr("charset")
		return !exists
	})

	if selection != nil {
		cs, _ := selection.Attr("charset")
		return NormaliseCharset(cs)
	}

	return ""
}

// Crawl fetches the HTML body and returns an Article
func (c Crawler) Crawl() *Article {

	article := new(Article)
	c.assignParseCandidate()
	c.assignHTML()

	if c.RawHTML == "" {
		return article
	}

	c.RawHTML = c.addSpacesBetweenTags(c.RawHTML)

	reader := strings.NewReader(c.RawHTML)
	document, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		panic(err.Error())
	}

	cs := c.GetCharset(document)
	if "" != cs && "UTF-8" != cs {
		c.RawHTML = UTF8encode(c.RawHTML, cs)
		reader = strings.NewReader(c.RawHTML)
		document, err = goquery.NewDocumentFromReader(reader)

		if nil != err {
			panic(err.Error())
		}
	}

	extractor := NewExtractor(c.config)

	startTime := time.Now().UnixNano()

	article.RawHTML, _ = document.Html()
	article.FinalURL = c.helper.url
	article.LinkHash = c.helper.linkHash
	article.Doc = document
	article.Title = extractor.GetTitle(document)
	article.MetaLang = extractor.GetMetaLanguage(document)
	article.MetaFavicon = extractor.GetFavicon(document)

	article.MetaDescription = extractor.GetMetaContentWithSelector(document, "meta[name#=(?i)^description$]")
	article.MetaKeywords = extractor.GetMetaContentWithSelector(document, "meta[name#=(?i)^keywords$]")
	article.CanonicalLink = extractor.GetCanonicalLink(document)
	if "" == article.CanonicalLink {
		article.CanonicalLink = article.FinalURL
	}
	article.Domain = extractor.GetDomain(article.CanonicalLink)
	article.Tags = extractor.GetTags(document)

	cleaner := NewCleaner(c.config)
	article.Doc = cleaner.clean(article.Doc)

	article.TopImage = OpenGraphResolver(article)
	if article.TopImage == "" {
		article.TopImage = WebPageResolver(article)
	}
	article.TopNode = extractor.calculateBestNode(article)
	if article.TopNode != nil {
		article.TopNode = extractor.postCleanup(article.TopNode)

		outputFormatter := new(outputFormatter)
		article.CleanedText, article.Links = outputFormatter.getFormattedText(article)

		videoExtractor := NewVideoExtractor()
		article.Movies = videoExtractor.GetVideos(article)
	}

	article.Delta = time.Now().UnixNano() - startTime

	return article
}

// In many cases, like at the end of each <li> element or between </span><span> tags,
// we need to add spaces or the text on either side will get joined together into one word.
// Also, add newlines after each </p> tag to preserve paragraphs.
func (c Crawler) addSpacesBetweenTags(text string) string {
	text = strings.Replace(text, "><", "> <", -1)
	text = strings.Replace(text, "</blockquote>", "</blockquote>\n", -1)
	text = strings.Replace(text, "<img ", "\n<img ", -1)
	return strings.Replace(text, "</p>", "</p>\n", -1)
}

func (c *Crawler) assignParseCandidate() {
	if c.RawHTML != "" {
		c.helper = NewRawHelper(c.url, c.RawHTML)
	} else {
		c.helper = NewURLHelper(c.url)
	}
}

func (c *Crawler) assignHTML() {
	if c.RawHTML == "" {
		cookieJar, _ := cookiejar.New(nil)
		client := &http.Client{
			Jar:     cookieJar,
			Timeout: c.config.timeout,
		}
		req, err := http.NewRequest("GET", c.url, nil)
		if err != nil {
			log.Println(err.Error())
			return
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_6_7) AppleWebKit/534.30 (KHTML, like Gecko) Chrome/12.0.742.91 Safari/534.30")
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err.Error())
			return
		}
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			c.RawHTML = string(contents)
		} else {
			log.Println(err.Error())
		}
	}
}
