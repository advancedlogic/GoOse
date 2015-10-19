package goose

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/advancedlogic/goquery"
	"github.com/rogpeppe/go-charset/charset"
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

// Crawl fetches the HTML body and returns an Article
func (c Crawler) Crawl() *Article {

	article := new(Article)
	c.assignParseCandidate()
	c.assignHTML()

	if c.RawHTML == "" {
		return article
	}

	reader := strings.NewReader(c.RawHTML)
	document, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		panic(err.Error())
	}

	attr := ""
	selection := document.Find("meta").EachWithBreak(func(i int, s *goquery.Selection) bool {
		attr, exists := s.Attr("http-equiv")
		if exists && attr == "Content-Type" {
			return false
		}
		return true
	})

	if selection != nil {
		attr, _ = selection.Attr("content")
		attr = strings.Replace(attr, " ", "", -1)

		if strings.HasPrefix(attr, "text/html;charset=") {
			cs := strings.TrimPrefix(attr, "text/html;charset=")
			cs = strings.ToLower(cs)

			if cs != "utf-8" {
				r, err := charset.NewReader(cs, strings.NewReader(c.RawHTML))
				if err != nil {
					// On error, skip the read
					c.RawHTML = ""
				} else {
					utf8, _ := ioutil.ReadAll(r)
					c.RawHTML = string(utf8)
				}
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
			article.CleanedText = outputFormatter.getFormattedText(article)

			videoExtractor := NewVideoExtractor()
			article.Movies = videoExtractor.GetVideos(article)
		}

		article.Delta = time.Now().UnixNano() - startTime

	} else {
		panic(err.Error())
	}
	return article
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
		if err == nil {
			req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_6_7) AppleWebKit/534.30 (KHTML, like Gecko) Chrome/12.0.742.91 Safari/534.30")
			resp, err := client.Do(req)
			if err == nil {
				defer resp.Body.Close()
				contents, err := ioutil.ReadAll(resp.Body)
				if err == nil {
					c.RawHTML = string(contents)
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
