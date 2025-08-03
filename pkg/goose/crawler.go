package goose

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

// HtmlRequester is a simple HTTP client for fetching web pages
type HtmlRequester struct {
	config Configuration
}

// NewHtmlRequester creates a new HTML requester
func NewHtmlRequester(config Configuration) HtmlRequester {
	return HtmlRequester{
		config: config,
	}
}

// fetchHTML fetches the HTML content from a URL
func (hr HtmlRequester) fetchHTML(targetURL string) (string, error) {
	// Parse URL
	_, err := url.Parse(targetURL)
	if err != nil {
		return "", err
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: hr.config.Timeout,
	}

	// Create request
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return "", err
	}

	// Set user agent
	req.Header.Set("User-Agent", hr.config.BrowserUserAgent)

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// SimpleCrawler is a basic crawler that extracts text
type SimpleCrawler struct {
	config Configuration
}

// NewCrawler creates a new crawler
func NewCrawler(config Configuration) SimpleCrawler {
	return SimpleCrawler{
		config: config,
	}
}

// Crawl extracts article content from HTML
func (c SimpleCrawler) Crawl(rawHTML string, targetURL string) (*Article, error) {
	article := &Article{
		RawHTML:  rawHTML,
		FinalURL: targetURL,
		Delta:    time.Now().Unix(),
	}

	// Simple HTML parsing - extract title and basic content
	// This is a very basic implementation
	article.Title = extractTitle(rawHTML)
	article.CleanedText = extractText(rawHTML)

	// Extract domain
	if u, err := url.Parse(targetURL); err == nil {
		article.Domain = u.Host
	}

	return article, nil
}

// extractTitle extracts the title from HTML (simple implementation)
func extractTitle(html string) string {
	// Very basic title extraction
	start := indexOf(html, "<title>")
	if start == -1 {
		start = indexOf(html, "<title ")
		if start == -1 {
			return ""
		}
		start = indexOf(html[start:], ">") + start + 1
	} else {
		start += 7
	}

	end := indexOf(html[start:], "</title>")
	if end == -1 {
		return ""
	}

	return html[start : start+end]
}

// extractText extracts basic text content (very simple implementation)
func extractText(html string) string {
	// This is a very basic text extraction
	// Remove script and style tags first
	text := removeTagContent(html, "script")
	text = removeTagContent(text, "style")

	// Remove HTML tags (very basic)
	result := ""
	inTag := false
	for _, char := range text {
		if char == '<' {
			inTag = true
		} else if char == '>' {
			inTag = false
		} else if !inTag {
			result += string(char)
		}
	}

	return result
}

// Helper functions
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func removeTagContent(html, tag string) string {
	openTag := "<" + tag
	closeTag := "</" + tag + ">"

	result := html
	for {
		start := indexOf(result, openTag)
		if start == -1 {
			break
		}

		// Find the end of the opening tag
		tagEnd := indexOf(result[start:], ">")
		if tagEnd == -1 {
			break
		}
		tagEnd += start + 1

		// Find the closing tag
		end := indexOf(result[tagEnd:], closeTag)
		if end == -1 {
			break
		}
		end += tagEnd + len(closeTag)

		// Remove the entire tag and its content
		result = result[:start] + result[end:]
	}

	return result
}
