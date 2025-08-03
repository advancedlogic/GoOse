package extractor

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/advancedlogic/GoOse/pkg/goose"
	"golang.org/x/net/html"
)

var normalizeWhitespaceRegexp = regexp.MustCompile(`[ \r\f\v\t]+`)
var normalizeNl = regexp.MustCompile(`\n{3,}`)
var multipleSpaces = regexp.MustCompile(`\n\s*\n\s*\n`)
var validURLRegex = regexp.MustCompile("^http[s]?://")

type outputFormatter struct {
	topNode  *goquery.Selection
	config   goose.Configuration
	language string
}

func (formatter *outputFormatter) getLanguage(lang string) string {
	if formatter.config.UseMetaLanguage && "" != lang {
		return lang
	}
	return formatter.config.TargetLanguage
}

func (formatter *outputFormatter) getTopNode() *goquery.Selection {
	return formatter.topNode
}

func (formatter *outputFormatter) getFormattedText(topNode *goquery.Selection, lang string) (output string, links []string) {
	formatter.topNode = topNode
	formatter.language = formatter.getLanguage(lang)
	if formatter.language == "" {
		formatter.language = formatter.config.TargetLanguage
	}
	formatter.removeNegativescoresNodes()
	links = formatter.linksToText()
	formatter.replaceTagsWithText()
	formatter.removeParagraphsWithFewWords()

	output = formatter.getOutputText()
	return output, links
}

func (formatter *outputFormatter) convertToText() string {
	var txts []string
	selections := formatter.topNode
	selections.Each(func(i int, s *goquery.Selection) {
		txt := s.Text()
		if txt != "" {
			// txt = txt //unescape
			txtLis := strings.Trim(txt, "\n")
			txts = append(txts, txtLis)
		}
	})
	return strings.Join(txts, "\n\n")
}

// check if this is a valid URL
func isValidURL(u string) bool {
	return validURLRegex.MatchString(u)
}

func (formatter *outputFormatter) linksToText() []string {
	var urlList []string
	links := formatter.topNode.Find("a")
	links.Each(func(i int, a *goquery.Selection) {
		imgs := a.Find("img")
		// ignore linked images
		if imgs.Length() == 0 {
			// save a list of URLs
			url, _ := a.Attr("href")
			if isValidURL(url) {
				urlList = append(urlList, url)
			}
			// replace <a> tag with its text contents
			replaceTagWithContents(a, whitelistedExtAtomTypes)

			// see whether we can collapse the parent node now
			replaceTagWithContents(a.Parent(), whitelistedTextAtomTypes)
		}
	})

	return urlList
}

// Text gets the combined text contents of each element in the set of matched
// elements, including their descendants.
//
// @see https://github.com/PuerkitoBio/goquery/blob/master/property.go
func (formatter *outputFormatter) Text(s *goquery.Selection) string {
	var buf bytes.Buffer

	// Slightly optimized vs calling Each: no single selection object created
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode && 0 == n.DataAtom { // NB: had to add the DataAtom check to avoid printing text twice when a textual node embeds another textual node
			// Keep newlines and spaces, like jQuery
			buf.WriteString(n.Data)
		}
		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	for _, n := range s.Nodes {
		f(n)
	}

	return buf.String()
}

func (formatter *outputFormatter) getOutputText() string {
	//out := formatter.topNode.Text()
	out := formatter.Text(formatter.topNode)
	out = normalizeWhitespaceRegexp.ReplaceAllString(out, " ")

	strArr := strings.Split(out, "\n")
	resArr := []string{}
	lastWasEmpty := false

	for _, v := range strArr {
		v = strings.TrimSpace(v)
		if v != "" {
			resArr = append(resArr, v)
			lastWasEmpty = false
		} else if !lastWasEmpty && len(resArr) > 0 {
			// Only add one empty line between paragraphs
			resArr = append(resArr, "")
			lastWasEmpty = true
		}
	}

	out = strings.Join(resArr, "\n")
	// More aggressive whitespace cleanup
	out = normalizeNl.ReplaceAllString(out, "\n\n")
	out = multipleSpaces.ReplaceAllString(out, "\n\n")
	
	// Final cleanup: remove leading/trailing whitespace
	out = strings.TrimSpace(out)
	
	// Additional cleanup for cases where content extraction has issues
	lines := strings.Split(out, "\n")
	cleanedLines := []string{}
	seenContent := make(map[string]bool)
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Skip lines that are clearly navigation or metadata
		if line != "" && !seenContent[line] && !formatter.isNavigationLine(line) {
			cleanedLines = append(cleanedLines, line)
			seenContent[line] = true
		} else if line == "" && len(cleanedLines) > 0 && cleanedLines[len(cleanedLines)-1] != "" {
			// Only add empty lines between different content blocks
			cleanedLines = append(cleanedLines, "")
		}
	}
	
	out = strings.Join(cleanedLines, "\n")
	out = strings.TrimSpace(out)
	
	return out
}

func (formatter *outputFormatter) removeNegativescoresNodes() {
	gravityItems := formatter.topNode.Find("*[gravityScore]")
	gravityItems.Each(func(i int, s *goquery.Selection) {
		var score int
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

func (formatter *outputFormatter) replaceTagsWithText() {
	for _, tag := range []string{"em", "strong", "b", "i", "span", "h1", "h2", "h3", "h4"} {
		nodes := formatter.topNode.Find(tag)
		nodes.Each(func(i int, node *goquery.Selection) {
			replaceTagWithContents(node, whitelistedTextAtomTypes)
		})
	}
}

func (formatter *outputFormatter) removeParagraphsWithFewWords() {
	language := formatter.language
	if language == "" {
		language = "en"
	}
	allNodes := formatter.topNode.Children()
	allNodes.Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		wordCount := len(strings.Fields(text))
		if wordCount < 5 && s.Find("object").Length() == 0 && s.Find("em").Length() == 0 {
			node := s.Get(0)
			if node.Parent != nil {
				node.Parent.RemoveChild(node)
			}
		}
	})
}

// isNavigationLine checks if a line of text is likely navigation or metadata
func (formatter *outputFormatter) isNavigationLine(line string) bool {
	if len(line) == 0 {
		return false
	}
	
	lowerLine := strings.ToLower(line)
	
	// Exact matches for navigation elements
	navExact := []string{
		"ad feedback", "cnn values your feedback", "how relevant is this ad to you?",
		"did you encounter any technical issues?", "video player was slow to load content",
		"video content never loaded", "ad froze or did not finish loading",
		"video content did not start after ad", "audio on ad was too loud",
		"other issues", "ad never loaded", "ad prevented/slowed the page from loading",
		"content moved around while ad loaded", "ad was repetitive to ads i've seen previously",
		"cancel", "submit", "thank you!", "your effort and contribution in providing this feedback is much appreciated.",
		"close", "close icon", "politics", "trump", "facts first", "cnn polls", "2025 elections",
		"more", "watch", "listen", "live tv", "subscribe", "sign in", "my account",
		"settings", "newsletters", "topics you follow", "sign out", "your cnn account",
		"sign in to your cnn account", "edition", "us", "international", "arabic", "español",
		"follow cnn politics", "crime + justice", "world", "africa", "americas", "asia",
		"australia", "china", "europe", "india", "middle east", "united kingdom",
		"business", "tech", "media", "calculators", "videos", "markets", "pre-markets",
		"after-hours", "fear & greed", "investing", "markets now", "nightcap", "health",
		"life, but better", "fitness", "food", "sleep", "mindfulness", "relationships",
		"cnn underscored", "electronics", "fashion", "beauty", "health & fitness",
		"home", "reviews", "deals", "gifts", "travel", "outdoors", "pets",
		"entertainment", "movies", "television", "celebrity", "innovate",
		"foreseeable future", "mission: ahead", "work transformed", "innovative cities",
		"style", "arts", "design", "architecture", "luxury", "video", "destinations",
		"food & drink", "stay", "sports", "pro football", "college football",
		"basketball", "baseball", "soccer", "olympics", "hockey", "science", "space",
		"life", "unearthed", "climate", "solutions", "weather", "ukraine-russia war",
		"israel-hamas war", "cnn headlines", "cnn shorts", "shows a-z", "cnn10",
		"cnn max", "cnn tv schedules", "flashdocs", "cnn 5 things",
		"chasing life with dr. sanjay gupta", "the assignment with audie cornish",
		"one thing", "tug of war", "cnn political briefing", "the axe files",
		"all there is with anderson cooper", "all cnn audio podcasts", "games",
		"daily crossword", "jumble crossword", "photo shuffle", "sudoblock", "sudoku",
		"5 things quiz", "about cnn", "photos", "investigations", "cnn profiles",
		"cnn leadership", "cnn newsletters", "work for cnn", "news", "terms of use",
		"privacy policy", "ad choices", "accessibility & cc", "about", "transcripts",
		"help center", "© 2025 cable news network. a warner bros. discovery company. all rights reserved.",
		"cnn sans ™ & © 2016 cable news network.", "facebook", "tweet", "email",
		"link", "link copied!", "follow", "see all topics", "donald trump",
	}
	
	for _, exact := range navExact {
		if lowerLine == exact {
			return true
		}
	}
	
	// Pattern matches
	navPatterns := []string{
		"min read", "updated", "published", "analysis by", "getty images",
		"reuters", "bloomberg", "afp", "via getty images", "jim watson/afp/getty images",
		"annabelle gordon/reuters", "jamie kelter davis/bloomberg/getty images",
	}
	
	for _, pattern := range navPatterns {
		if strings.Contains(lowerLine, pattern) {
			return true
		}
	}
	
	// Check for very short lines that are likely navigation
	if len(strings.TrimSpace(line)) < 3 {
		return true
	}
	
	// Check for lines that are just numbers or punctuation
	if strings.TrimSpace(line) == "•" || strings.TrimSpace(line) == "1." || strings.TrimSpace(line) == "2." {
		return true
	}
	
	return false
}
