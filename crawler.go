package goose

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

// Crawler can fetch the target HTML page
type Crawler struct {
	config  Configuration
	url     string
	RawHTML string
	helper  Helper
}

// NewCrawler returns a crawler object initialised with the URL and the [optional] raw HTML body
func NewCrawler(config Configuration, url string, RawHTML string) Crawler {
	return Crawler{
		config:  config,
		url:     url,
		RawHTML: RawHTML,
	}
}

// convert to UTF-8, skipping invalid byte sequences
// @see http://stackoverflow.com/questions/32512500/ignore-illegal-bytes-when-decoding-text-with-go
func utf8encode(raw string, sourceCharset string) string {
	enc, _ := charset.Lookup(sourceCharset)
	dst := make([]byte, len(raw))
	d := enc.NewDecoder()

	var (
		in  int
		out int
	)
	for in < len(raw) {
		// Do the transformation
		ndst, nsrc, err := d.Transform(dst[out:], []byte(raw[in:]), true)
		in += nsrc
		out += ndst
		if err == nil {
			// Completed transformation
			break
		}
		if err == transform.ErrShortDst {
			// Our output buffer is too small, so we need to grow it
			t := make([]byte, (cap(dst)+1)*2)
			copy(t, dst)
			dst = t
			continue
		}
		// We're here because of at least one illegal character. Skip over the current rune
		// and try again.
		_, width := utf8.DecodeRuneInString(raw[in:])
		in += width
	}
	return string(dst)
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
			cs := strings.ToLower(strings.TrimPrefix(attr, "text/html;charset="))
			if cs != "utf-8" {
				c.RawHTML = utf8encode(c.RawHTML, cs)
				reader = strings.NewReader(c.RawHTML)
				document, err = goquery.NewDocumentFromReader(reader)
			}
		}
	}

	if err == nil {
		extractor := NewExtractor(c.config)
		html, _ := document.Html()

		startTime := time.Now().UnixNano()
		article.RawHTML = html
		article.FinalURL = c.helper.url
		article.LinkHash = c.helper.linkHash
		article.Doc = document
		article.Title = extractor.getTitle(article)
		article.MetaLang = extractor.getMetaLanguage(article)
		article.MetaFavicon = extractor.getFavicon(article)

		article.MetaDescription = extractor.getMetaContentWithSelector(article, "meta[name#=(?i)^description$]")
		article.MetaKeywords = extractor.getMetaContentWithSelector(article, "meta[name#=(?i)^keywords$]")
		article.CanonicalLink = extractor.getCanonicalLink(article)
		article.Domain = extractor.getDomain(article)
		article.Tags = extractor.getTags(article)

		cleaner := NewCleaner(c.config)
		article.Doc = cleaner.clean(article)

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

	} else {
		panic(err.Error())
	}
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
