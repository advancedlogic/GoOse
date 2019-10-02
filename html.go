package goose

import (
	resty "github.com/go-resty/resty/v2"
)

type HtmlRequester interface {
	fetchHTML(string) (string, error)
}

// Crawler can fetch the target HTML page
type htmlrequester struct {
	config Configuration
}

// NewCrawler returns a crawler object initialised with the URL and the [optional] raw HTML body
func NewHtmlRequester(config Configuration) HtmlRequester {
	return htmlrequester{
		config: config,
	}
}

func (hr htmlrequester) fetchHTML(url string) (string, error) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		//TODO:Set timeout
		// req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_6_7) AppleWebKit/534.30 (KHTML, like Gecko) Chrome/12.0.742.91 Safari/534.30")
		Get(url)

	// 	contents, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return "", err
	// }
	// c.RawHTML = string(contents)

	// if err = resp.Body.Close(); err != nil {
	// 	return "", err
	// }

	if err != nil {
		return "", err
	}
	if resp.IsError() {
		// TODO: real error return
		return "", nil
	}
	return resp.String(), nil
}
