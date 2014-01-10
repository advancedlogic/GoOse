package goose

func init() {
	LoadGooseStopwords("resources/text")
}

type Goose struct {
	config Configuration
}

func NewGoose() Goose {

	return Goose{
		config: GetDefualtConfiguration(),
	}
}

func GetGoosFromConfiguration(config Configuration) Goose {
	return Goose{
		config: config,
	}
}

func (gs Goose) ExtractFromUrl(url string) *Article {
	cc := NewCrawlCandidate(url, "")
	cc.config = gs.config
	return cc.Crawl()
}

func (gs Goose) ExtractFromRawHtml(url string, rawHtml string) *Article {
	cc := NewCrawlCandidate(url, rawHtml)
	cc.config = gs.config
	return cc.Crawl()
}
