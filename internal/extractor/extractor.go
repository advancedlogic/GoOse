package extractor

import (
	"container/list"
	"log"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/araddon/dateparse"
	"github.com/fatih/set"
	"github.com/gigawattio/window"
	"github.com/jaytaylor/html2text"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/advancedlogic/GoOse/internal/utils"
	"github.com/advancedlogic/GoOse/pkg/goose"
)

const defaultLanguage = "en"

var motleyReplacement = "&#65533;" // U+FFFD (decimal 65533) is the "replacement character".
//var escapedFragmentReplacement = regexp.MustCompile("#!")
//var titleReplacements = regexp.MustCompile("&raquo;")

var titleDelimiters = []string{
	"|",
	" - ",
	" — ",
	"»",
	":",
}

var aRelTagSelector = "a[rel=tag]"
var aHrefTagSelector = [...]string{"/tag/", "/tags/", "/topic/", "?keyword"}

//var langRegEx = "^[A-Za-z]{2}$"

// ContentExtractor can parse the HTML and fetch various properties
type ContentExtractor struct {
	config goose.Configuration
}

// NewExtractor returns a configured HTML parser
func NewExtractor(config goose.Configuration) ContentExtractor {
	return ContentExtractor{
		config: config,
	}
}

// if the article has a title set in the source, use that
func (extr *ContentExtractor) getTitleUnmodified(document *goquery.Document) string {
	title := ""

	titleElement := document.Find("title")
	if titleElement != nil && titleElement.Size() > 0 {
		title = titleElement.Text()
	}

	if title == "" {
		ogTitleElement := document.Find(`meta[property="og:title"]`)
		if ogTitleElement != nil && ogTitleElement.Size() > 0 {
			title, _ = ogTitleElement.Attr("content")
		}
	}

	if title == "" {
		titleElement = document.Find("post-title,headline")
		if titleElement == nil || titleElement.Size() == 0 {
			return title
		}
		title = titleElement.Text()
	}
	return title
}

// GetTitleFromUnmodifiedTitle returns the title from the unmodified one
func (extr *ContentExtractor) GetTitleFromUnmodifiedTitle(title string) string {
	originalTitle := title
	for _, delimiter := range titleDelimiters {
		if strings.Contains(title, delimiter) {
			parts := strings.Split(title, delimiter)
			if extr.config.Debug {
				log.Printf("Found delimiter '%s', split into %d parts\n", delimiter, len(parts))
				for i, part := range parts {
					log.Printf("  Part %d: %q (len=%d)\n", i, part, len(part))
				}
			}
			title = extr.splitTitle(parts)
			if extr.config.Debug {
				log.Printf("After splitTitle: %q\n", title)
			}
			break
		}
	}

	title = strings.Replace(title, motleyReplacement, "", -1)

	if extr.config.Debug {
		log.Printf("Original title: %q\n", originalTitle)
		log.Printf("Final title: %q\n", title)
	}

	return strings.TrimSpace(title)
}

// GetTitle returns the title set in the source, if the article has one
func (extr *ContentExtractor) GetTitle(document *goquery.Document) string {
	title := extr.getTitleUnmodified(document)
	if extr.config.Debug {
		log.Printf("Unmodified title: %q\n", title)
	}
	return extr.GetTitleFromUnmodifiedTitle(title)
}

func (extr *ContentExtractor) splitTitle(titles []string) string {
	// For common patterns like "Article Title - Site Name", prefer the first part
	if len(titles) >= 2 {
		// Trim spaces from all parts
		for i := range titles {
			titles[i] = strings.TrimSpace(titles[i])
		}
		
		// Check if last part looks like a site name (common pattern)
		lastPart := titles[len(titles)-1]
		// Common site name patterns
		if len(titles) == 2 && (strings.Contains(lastPart, "News") || 
			strings.Contains(lastPart, "BBC") || 
			strings.Contains(lastPart, "CNN") || 
			strings.Contains(lastPart, "ABC") ||
			strings.Contains(lastPart, "Times") || 
			strings.Contains(lastPart, "Post") || 
			strings.Contains(lastPart, "Journal") ||
			len(lastPart) < 20) {
			// Return the first part
			title := strings.Replace(titles[0], "&raquo;", "»", -1)
			return title
		}
	}
	
	// Fallback to the original logic - choose the longest part
	largeTextLength := 0
	largeTextIndex := 0
	for i, current := range titles {
		if len(current) > largeTextLength {
			largeTextLength = len(current)
			largeTextIndex = i
		}
	}
	title := titles[largeTextIndex]
	title = strings.Replace(title, "&raquo;", "»", -1)
	return title
}

// GetMetaLanguage returns the meta language set in the source, if the article has one
func (extr *ContentExtractor) GetMetaLanguage(document *goquery.Document) string {
	var language string
	shtml := document.Find("html")
	attr, _ := shtml.Attr("lang")
	if attr == "" {
		attr, _ = document.Attr("lang")
	}
	if attr == "" {
		selection := document.Find("meta").EachWithBreak(func(i int, s *goquery.Selection) bool {
			var exists bool
			attr, exists = s.Attr("http-equiv")
			if exists && attr == "content-language" {
				return false
			}
			return true
		})
		if selection != nil {
			attr, _ = selection.Attr("content")
		}
	}
	idx := strings.LastIndex(attr, "-")
	if idx == -1 {
		language = attr
	} else {
		language = attr[0:idx]
	}

	_, ok := utils.Sw[language]

	if language == "" || !ok {
		language = extr.config.StopWords.SimpleLanguageDetector(shtml.Text())
		if language == "" {
			language = defaultLanguage
		}
	}

	extr.config.TargetLanguage = language
	return language
}

// GetFavicon returns the favicon set in the source, if the article has one
func (extr *ContentExtractor) GetFavicon(document *goquery.Document) string {
	favicon := ""
	document.Find("link").EachWithBreak(func(i int, s *goquery.Selection) bool {
		attr, exists := s.Attr("rel")
		if exists && strings.Contains(attr, "icon") {
			favicon, _ = s.Attr("href")
			return false
		}
		return true
	})
	return favicon
}

// GetMetaContentWithSelector returns the content attribute of meta tag matching the selector
func (extr *ContentExtractor) GetMetaContentWithSelector(document *goquery.Document, selector string) string {
	selection := document.Find(selector)
	content, _ := selection.Attr("content")
	return strings.TrimSpace(content)
}

// GetMetaContent returns the content attribute of meta tag with the given property name
func (extr *ContentExtractor) GetMetaContent(document *goquery.Document, metaName string) string {
	content := ""
	document.Find("meta").EachWithBreak(func(i int, s *goquery.Selection) bool {
		attr, exists := s.Attr("name")
		if exists && attr == metaName {
			content, _ = s.Attr("content")
			return false
		}
		attr, exists = s.Attr("itemprop")
		if exists && attr == metaName {
			content, _ = s.Attr("content")
			return false
		}
		return true
	})
	return content
}

// GetMetaContents returns all the meta tags as name->content pairs
func (extr *ContentExtractor) GetMetaContents(document *goquery.Document, metaNames *set.Set) map[string]string {
	contents := make(map[string]string)
	counter := metaNames.Size()
	document.Find("meta").EachWithBreak(func(i int, s *goquery.Selection) bool {
		attr, exists := s.Attr("name")
		if exists && metaNames.Has(attr) {
			content, _ := s.Attr("content")
			contents[attr] = content
			counter--
			if counter < 0 {
				return false
			}
		}
		return true
	})
	return contents
}

// GetMetaDescription returns the meta description set in the source, if the article has one
func (extr *ContentExtractor) GetMetaDescription(document *goquery.Document) string {
	return extr.GetMetaContent(document, "description")
}

// GetMetaKeywords returns the meta keywords set in the source, if the article has them
func (extr *ContentExtractor) GetMetaKeywords(document *goquery.Document) string {
	return extr.GetMetaContent(document, "keywords")
}

// GetMetaAuthor returns the meta author set in the source, if the article has one
func (extr *ContentExtractor) GetMetaAuthor(document *goquery.Document) string {
	return extr.GetMetaContent(document, "author")
}

// GetMetaContentLocation returns the meta content location set in the source, if the article has one
func (extr *ContentExtractor) GetMetaContentLocation(document *goquery.Document) string {
	return extr.GetMetaContent(document, "contentLocation")
}

// GetCanonicalLink returns the meta canonical link set in the source
func (extr *ContentExtractor) GetCanonicalLink(document *goquery.Document) string {
	metas := document.Find("link[rel=canonical]")
	if metas.Length() > 0 {
		meta := metas.First()
		href, _ := meta.Attr("href")
		href = strings.Trim(href, "\n")
		href = strings.Trim(href, " ")
		if href != "" {
			return href
		}
	}
	return ""
}

// GetDomain extracts the domain from a link
func (extr *ContentExtractor) GetDomain(canonicalLink string) string {
	u, err := url.Parse(canonicalLink)
	if err == nil {
		return u.Host
	}
	return ""
}

// GetTags returns the tags set in the source, if the article has them
func (extr *ContentExtractor) GetTags(document *goquery.Document) *set.Set {
	tags := set.New(set.ThreadSafe).(*set.Set)
	selections := document.Find(aRelTagSelector)
	selections.Each(func(i int, s *goquery.Selection) {
		tags.Add(s.Text())
	})
	selections = document.Find("a")
	selections.Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			for _, part := range aHrefTagSelector {
				if strings.Contains(href, part) {
					tags.Add(s.Text())
				}
			}
		}
	})

	return tags
}

// GetPublishDate returns the publication date, if one can be located.
func (extr *ContentExtractor) GetPublishDate(document *goquery.Document) *time.Time {
	raw, err := document.Html()
	if err != nil {
		log.Printf("Error converting document HTML nodes to raw HTML: %s (publish date detection aborted)\n", err)
		return nil
	}

	text, err := html2text.FromString(raw)
	if err != nil {
		log.Printf("Error converting document HTML to plaintext: %s (publish date detection aborted)\n", err)
		return nil
	}

	text = strings.ToLower(text)

	// Simplify months because the dateparse pkg only handles abbreviated.
	for k, v := range map[string]string{
		"january":  "jan",
		"march":    "mar",
		"february": "feb",
		"april":    "apr",
		// "may":       "may", // Pointless.
		"june":      "jun",
		"august":    "aug",
		"september": "sep",
		"sept":      "sep",
		"october":   "oct",
		"november":  "nov",
		"december":  "dec",
		"th,":       ",", // Strip day number suffixes.
		"rd,":       ",",
	} {
		text = strings.Replace(text, k, v, -1)
	}
	text = strings.Replace(text, "\n", " ", -1)
	text = regexp.MustCompile(" +").ReplaceAllString(text, " ")

	tuple1 := strings.Split(text, " ")

	var (
		expr  = regexp.MustCompile("[0-9]")
		ts    time.Time
		found bool
	)
	for _, n := range []int{3, 4, 5, 2, 6} {
		for _, win := range window.Rolling(tuple1, n) {
			if !expr.MatchString(strings.Join(win, " ")) {
				continue
			}

			input := strings.Join(win, " ")
			ts, err = dateparse.ParseAny(input)
			if err == nil && ts.Year() > 0 && ts.Month() > 0 && ts.Day() > 0 {
				found = true
				break
			}

			// Try injecting a comma for dateparse.
			win[1] = win[1] + ","
			input = strings.Join(win, " ")
			ts, err = dateparse.ParseAny(input)
			if err == nil && ts.Year() > 0 && ts.Month() > 0 && ts.Day() > 0 {
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if found {
		return &ts
	}
	return nil
}

// GetCleanTextAndLinks parses the main HTML node for text and links
func (extr *ContentExtractor) GetCleanTextAndLinks(topNode *goquery.Selection, lang string) (string, []string) {
	outputFormatter := new(outputFormatter)
	outputFormatter.config = extr.config
	return outputFormatter.getFormattedText(topNode, lang)
}

// CalculateBestNode checks for the HTML node most likely to contain the main content.
// we're going to start looking for where the clusters of paragraphs are. We'll score a cluster based on the number of stopwords
// and the number of consecutive paragraphs together, which should form the cluster of text that this node is around
// also store on how high up the paragraphs are, comments are usually at the bottom and should get a lower score
func (extr *ContentExtractor) CalculateBestNode(document *goquery.Document) *goquery.Selection {
	// First try site-specific selectors for known news sites
	if siteSpecificNode := extr.tryNewsSelectors(document); siteSpecificNode != nil {
		return siteSpecificNode
	}
	
	var topNode *goquery.Selection
	nodesToCheck := extr.nodesToCheck(document)
	if extr.config.Debug {
		log.Printf("Nodes to check %d\n", len(nodesToCheck))
	}
	startingBoost := 1.0
	cnt := 0
	i := 0
	parentNodes := set.New(set.ThreadSafe).(*set.Set)
	nodesWithText := list.New()
	for _, node := range nodesToCheck {
		textNode := node.Text()
		ws := extr.config.StopWords.StopWordsCount(extr.config.TargetLanguage, textNode)
		highLinkDensity := extr.isHighLinkDensity(node)
		
		// Boost scoring for nodes that look like article content
		articleBoost := extr.getArticleContentBoost(node)
		adjustedWs := ws + articleBoost
		
		if adjustedWs > 2 && !highLinkDensity {
			nodesWithText.PushBack(node)
		}
	}
	nodesNumber := nodesWithText.Len()
	negativeScoring := 0
	bottomNegativeScoring := float64(nodesNumber) * 0.25

	if extr.config.Debug {
		log.Printf("About to inspect num of nodes with text %d\n", nodesNumber)
	}

	for n := nodesWithText.Front(); n != nil; n = n.Next() {
		node := n.Value.(*goquery.Selection)
		boostScore := 0.0
		if extr.isBoostable(node) {
			if cnt >= 0 {
				boostScore = float64((1.0 / startingBoost) * 50)
				startingBoost++
			}
		}

		if nodesNumber > 15 {
			if float64(nodesNumber-i) <= bottomNegativeScoring {
				booster := bottomNegativeScoring - float64(nodesNumber-i)
				boostScore = -math.Pow(booster, 2.0)
				negScore := math.Abs(boostScore) + float64(negativeScoring)
				if negScore > 40 {
					boostScore = 5.0
				}
			}
		}

		if extr.config.Debug {
			log.Printf("Location Boost Score %1.5f on iteration %d id='%s' class='%s'\n", boostScore, i, extr.config.Parser.Name("id", node), extr.config.Parser.Name("class", node))
		}
		textNode := node.Text()
		ws := extr.config.StopWords.StopWordsCount(extr.config.TargetLanguage, textNode)
		upScore := ws + int(boostScore)
		parentNode := node.Parent()
		extr.updateScore(parentNode, upScore)
		extr.updateNodeCount(parentNode, 1)
		if !parentNodes.Has(parentNode) {
			parentNodes.Add(parentNode)
		}
		parentParentNode := parentNode.Parent()
		if parentParentNode != nil {
			extr.updateNodeCount(parentParentNode, 1)
			extr.updateScore(parentParentNode, upScore/2.0)
			if !parentNodes.Has(parentParentNode) {
				parentNodes.Add(parentParentNode)
			}
		}
		cnt++
		i++
	}

	topNodeScore := 0
	parentNodesArray := parentNodes.List()
	for _, p := range parentNodesArray {
		e := p.(*goquery.Selection)
		if extr.config.Debug {
			log.Printf("ParentNode: score=%s nodeCount=%s id='%s' class='%s'\n", extr.config.Parser.Name("gravityScore", e), extr.config.Parser.Name("gravityNodes", e), extr.config.Parser.Name("id", e), extr.config.Parser.Name("class", e))
		}
		score := extr.getScore(e)
		if score >= topNodeScore {
			topNode = e
			topNodeScore = score
		}
		if topNode == nil {
			topNode = e
		}
	}
	return topNode
}

// returns the gravityScore as an integer from this node
func (extr *ContentExtractor) getScore(node *goquery.Selection) int {
	return extr.getNodeGravityScore(node)
}

func (extr *ContentExtractor) getNodeGravityScore(node *goquery.Selection) int {
	grvScoreString, exists := node.Attr("gravityScore")
	if !exists {
		return 0
	}
	grvScore, err := strconv.Atoi(grvScoreString)
	if err != nil {
		return 0
	}
	return grvScore
}

// adds a score to the gravityScore Attribute we put on divs
// we'll get the current score then add the score we're passing in to the current
func (extr *ContentExtractor) updateScore(node *goquery.Selection, addToScore int) {
	currentScore := 0
	var err error
	scoreString, _ := node.Attr("gravityScore")
	if scoreString != "" {
		currentScore, err = strconv.Atoi(scoreString)
		if err != nil {
			currentScore = 0
		}
	}
	newScore := currentScore + addToScore
	extr.config.Parser.SetAttr(node, "gravityScore", strconv.Itoa(newScore))
}

// stores how many decent nodes are under a parent node
func (extr *ContentExtractor) updateNodeCount(node *goquery.Selection, addToCount int) {
	currentScore := 0
	var err error
	scoreString, _ := node.Attr("gravityNodes")
	if scoreString != "" {
		currentScore, err = strconv.Atoi(scoreString)
		if err != nil {
			currentScore = 0
		}
	}
	newScore := currentScore + addToCount
	extr.config.Parser.SetAttr(node, "gravityNodes", strconv.Itoa(newScore))
}

// a lot of times the first paragraph might be the caption under an image so we'll want to make sure if we're going to
// boost a parent node that it should be connected to other paragraphs, at least for the first n paragraphs
// so we'll want to make sure that the next sibling is a paragraph and has at least some substantial weight to it
func (extr *ContentExtractor) isBoostable(node *goquery.Selection) bool {
	stepsAway := 0
	next := node.Next()
	for next != nil && stepsAway < node.Siblings().Length() {
		currentNodeTag := node.Get(0).DataAtom.String()
		if currentNodeTag == "p" {
			if stepsAway >= 3 {
				if extr.config.Debug {
					log.Println("Next paragraph is too far away, not boosting")
				}
				return false
			}

			paraText := node.Text()
			ws := extr.config.StopWords.StopWordsCount(extr.config.TargetLanguage, paraText)
			if ws > 5 {
				if extr.config.Debug {
					log.Println("We're gonna boost this node, seems content")
				}
				return true
			}
		}

		stepsAway++
		next = next.Next()
	}

	return false
}

// returns a list of nodes we want to search on like paragraphs and tables
func (extr *ContentExtractor) nodesToCheck(doc *goquery.Document) []*goquery.Selection {
	var output []*goquery.Selection
	tags := []string{"p", "pre", "td"}
	for _, tag := range tags {
		selections := doc.Children().Find(tag)
		if selections != nil {
			selections.Each(func(i int, s *goquery.Selection) {
				// Skip nodes that are clearly navigation or non-content
				if extr.isLikelyNonContent(s) {
					return
				}
				output = append(output, s)
			})
		}
	}
	return output
}

// Helper function to identify nodes that are likely non-content
func (extr *ContentExtractor) isLikelyNonContent(node *goquery.Selection) bool {
	// Check parent hierarchy for navigation indicators
	for parent := node.Parent(); parent != nil && parent.Length() > 0; parent = parent.Parent() {
		class, hasClass := parent.Attr("class")
		id, hasId := parent.Attr("id")
		
		if hasClass {
			class = strings.ToLower(class)
			if strings.Contains(class, "nav") || strings.Contains(class, "menu") || 
			   strings.Contains(class, "header") || strings.Contains(class, "footer") ||
			   strings.Contains(class, "sidebar") || strings.Contains(class, "aside") ||
			   strings.Contains(class, "ad") || strings.Contains(class, "banner") ||
			   strings.Contains(class, "breadcrumb") || strings.Contains(class, "related") {
				return true
			}
		}
		
		if hasId {
			id = strings.ToLower(id)
			if strings.Contains(id, "nav") || strings.Contains(id, "menu") || 
			   strings.Contains(id, "header") || strings.Contains(id, "footer") ||
			   strings.Contains(id, "sidebar") || strings.Contains(id, "aside") ||
			   strings.Contains(id, "ad") || strings.Contains(id, "banner") {
				return true
			}
		}
		
		// Check tag type
		tagName := parent.Get(0).DataAtom.String()
		if tagName == "nav" || tagName == "header" || tagName == "footer" || tagName == "aside" {
			return true
		}
	}
	
	// Check if the node itself has very short text (likely navigation link)
	text := strings.TrimSpace(node.Text())
	if len(text) < 10 {
		return true
	}
	
	return false
}

// getArticleContentBoost provides additional scoring for nodes that appear to be article content
func (extr *ContentExtractor) getArticleContentBoost(node *goquery.Selection) int {
	boost := 0
	
	// Check parent hierarchy for article-related classes/ids
	for parent := node.Parent(); parent != nil && parent.Length() > 0; parent = parent.Parent() {
		class, hasClass := parent.Attr("class")
		id, hasId := parent.Attr("id")
		
		if hasClass {
			class = strings.ToLower(class)
			if strings.Contains(class, "article") || strings.Contains(class, "content") ||
			   strings.Contains(class, "story") || strings.Contains(class, "post") ||
			   strings.Contains(class, "entry") || strings.Contains(class, "main") ||
			   strings.Contains(class, "body") || strings.Contains(class, "text") {
				boost += 10
			}
		}
		
		if hasId {
			id = strings.ToLower(id)
			if strings.Contains(id, "article") || strings.Contains(id, "content") ||
			   strings.Contains(id, "story") || strings.Contains(id, "post") ||
			   strings.Contains(id, "entry") || strings.Contains(id, "main") {
				boost += 10
			}
		}
		
		// Check for semantic HTML5 tags
		tagName := parent.Get(0).DataAtom.String()
		if tagName == "article" || tagName == "main" {
			boost += 15
		}
	}
	
	// Penalize nodes that seem to be in navigation or sidebars
	text := strings.TrimSpace(node.Text())
	if len(text) > 100 { // Long text is more likely to be content
		boost += 5
	}
	
	// Look for paragraph length - articles typically have substantial paragraphs
	if len(text) > 200 {
		boost += 5
	}
	
	return boost
}

// tryNewsSelectors attempts to find article content using site-specific selectors
func (extr *ContentExtractor) tryNewsSelectors(document *goquery.Document) *goquery.Selection {
	// Common article content selectors used by major news sites
	selectors := []string{
		".article__content", // CNN and other major news sites
		"article[role='main']",
		"[data-module='ArticleBody']",
		".article-body",
		".story-body",
		".post-content",
		".entry-content", 
		".content-body",
		"main article",
		"[role='main'] article",
		".article-wrap .article-content",
		".story-content",
		".post-body",
		"#article-body",
		"#story-body",
		"#content .article",
		".story__content",
		".post__content",
		"[data-testid='article-body']",
		"[data-testid='story-body']",
	}
	
	for _, selector := range selectors {
		selection := document.Find(selector)
		if selection.Length() > 0 {
			// Validate that this looks like article content
			text := strings.TrimSpace(selection.Text())
			if len(text) > 200 { // Must have substantial content
				// Check for reasonable paragraph count
				paragraphs := selection.Find("p")
				if paragraphs.Length() >= 3 { // Should have multiple paragraphs
					// Additional validation: ensure it's not mostly navigation
					if !extr.isHighLinkDensity(selection) && extr.hasGoodContentSignals(selection) {
						if extr.config.Debug {
							log.Printf("Found article content using selector: %s (text length: %d, paragraphs: %d)\n", 
								selector, len(text), paragraphs.Length())
						}
						// Extract only the paragraph content, not the entire container
						return extr.extractParagraphContent(selection)
					}
				}
			}
		}
	}
	
	// Try looking for elements with substantial text content that aren't navigation
	var bestCandidate *goquery.Selection
	var bestScore int
	
	document.Find("div, article, section").Each(func(i int, s *goquery.Selection) {
		class, _ := s.Attr("class")
		id, _ := s.Attr("id")
		
		// Look for likely content containers
		if strings.Contains(strings.ToLower(class), "content") || 
		   strings.Contains(strings.ToLower(class), "article") ||
		   strings.Contains(strings.ToLower(class), "story") ||
		   strings.Contains(strings.ToLower(id), "content") ||
		   strings.Contains(strings.ToLower(id), "article") ||
		   strings.Contains(strings.ToLower(id), "story") {
			
			text := strings.TrimSpace(s.Text())
			if len(text) > 500 { // Substantial content
				paragraphs := s.Find("p")
				if paragraphs.Length() >= 5 { // Multiple paragraphs
					// Check that it's not mostly links (navigation)
					if !extr.isHighLinkDensity(s) {
						score := len(text) + (paragraphs.Length() * 50)
						if score > bestScore {
							bestCandidate = s
							bestScore = score
							if extr.config.Debug {
								log.Printf("Found potential article content by class/id: %s %s (text length: %d, score: %d)\n", 
									class, id, len(text), score)
							}
						}
					}
				}
			}
		}
	})
	
	return bestCandidate
}

// hasGoodContentSignals checks if a node contains signals that indicate it's article content
func (extr *ContentExtractor) hasGoodContentSignals(node *goquery.Selection) bool {
	text := strings.TrimSpace(node.Text())
	
	// Check for article-like sentence structure (sentences ending with periods)
	sentences := strings.Split(text, ".")
	if len(sentences) < 3 {
		return false // Too few sentences for an article
	}
	
	// Check average sentence length (articles have substantial sentences)
	totalLength := 0
	validSentences := 0
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if len(sentence) > 20 { // Only count substantial sentences
			totalLength += len(sentence)
			validSentences++
		}
	}
	
	if validSentences < 3 {
		return false
	}
	
	avgSentenceLength := totalLength / validSentences
	if avgSentenceLength < 50 { // Articles typically have longer sentences
		return false
	}
	
	// Check for common navigation patterns to exclude
	lowerText := strings.ToLower(text)
	navigationWords := []string{
		"sign in", "sign out", "subscribe", "newsletter", "account",
		"home", "news", "sports", "weather", "politics", "business",
		"watch", "listen", "live tv", "more", "follow", "settings",
		"crime", "world", "health", "entertainment", "travel",
		"calculators", "markets", "investing", "fashion", "beauty",
		"games", "crossword", "photos", "investigations", "profiles",
	}
	
	navigationCount := 0
	for _, word := range navigationWords {
		if strings.Contains(lowerText, word) {
			navigationCount++
		}
	}
	
	// If it contains many navigation words, it's likely not article content
	if navigationCount > 5 {
		return false
	}
	
	return true
}

// extractParagraphContent creates a new selection containing only the substantial paragraphs from the article
func (extr *ContentExtractor) extractParagraphContent(selection *goquery.Selection) *goquery.Selection {
	// Create a new document fragment with only the article paragraphs
	paragraphs := selection.Find("p")
	var cleanParagraphs []*goquery.Selection
	
	paragraphs.Each(func(i int, p *goquery.Selection) {
		text := strings.TrimSpace(p.Text())
		// Only include paragraphs with substantial content
		if len(text) > 30 {
			// Skip paragraphs that look like metadata or navigation
			lowerText := strings.ToLower(text)
			if !strings.Contains(lowerText, "updated") && 
			   !strings.Contains(lowerText, "published") &&
			   !strings.Contains(lowerText, "min read") &&
			   !strings.Contains(lowerText, "follow") &&
			   !strings.Contains(lowerText, "subscribe") &&
			   !strings.Contains(lowerText, "sign in") &&
			   !strings.Contains(lowerText, "analysis by") &&
			   !strings.Contains(lowerText, "see all topics") {
				cleanParagraphs = append(cleanParagraphs, p)
			}
		}
	})
	
	// If we found good paragraphs, return the first one's parent and modify it
	if len(cleanParagraphs) > 0 {
		// Return the original selection but with filtered content
		// Remove all non-paragraph children first
		selection.Children().Each(func(i int, child *goquery.Selection) {
			if child.Get(0).DataAtom.String() != "p" {
				// Check if this is one of our clean paragraphs
				isCleanParagraph := false
				for _, cleanP := range cleanParagraphs {
					if child.Get(0) == cleanP.Get(0) {
						isCleanParagraph = true
						break
					}
				}
				if !isCleanParagraph {
					extr.config.Parser.RemoveNode(child)
				}
			}
		})
		return selection
	}
	
	return selection
}

// checks the density of links within a node, is there not much text and most of it contains bad links?
// if so it's no good
func (extr *ContentExtractor) isHighLinkDensity(node *goquery.Selection) bool {
	links := node.Find("a")
	if links == nil || links.Size() == 0 {
		return false
	}
	text := node.Text()
	words := strings.Split(text, " ")
	nwords := len(words)
	var sb []string
	links.Each(func(i int, s *goquery.Selection) {
		linkText := s.Text()
		sb = append(sb, linkText)
	})
	linkText := strings.Join(sb, "")
	linkWords := strings.Split(linkText, " ")
	nlinkWords := len(linkWords)
	nlinks := links.Size()
	
	// Avoid division by zero
	if nwords == 0 {
		return true
	}
	
	linkDivisor := float64(nlinkWords) / float64(nwords)
	score := linkDivisor * float64(nlinks)

	// More aggressive link density detection for better content isolation
	// Navigation menus typically have high link density
	if nlinks > 5 && linkDivisor > 0.3 {
		return true
	}
	
	// If more than 60% of words are in links, it's likely navigation
	if linkDivisor > 0.6 {
		return true
	}

	if extr.config.Debug {
		var logText string
		if len(node.Text()) >= 51 {
			logText = node.Text()[0:50]
		} else {
			logText = node.Text()
		}
		log.Printf("Calculated link density score as %1.5f for node %s (links: %d, linkDivisor: %1.3f)\n", score, logText, nlinks, linkDivisor)
	}
	if score > 0.8 {  // Lowered from 1.0 to be more aggressive
		return true
	}
	return false
}

func (extr *ContentExtractor) isTableAndNoParaExist(selection *goquery.Selection) bool {
	subParagraph := selection.Find("p")
	subParagraph.Each(func(i int, s *goquery.Selection) {
		txt := s.Text()
		if len(txt) < 25 {
			node := s.Get(0)
			parent := node.Parent
			parent.RemoveChild(node)
		}
	})

	subParagraph2 := selection.Find("p")
	if subParagraph2.Length() == 0 && selection.Get(0).DataAtom.String() != "td" {
		return true
	}
	return false
}

func (extr *ContentExtractor) isNodescoreThresholdMet(node *goquery.Selection, e *goquery.Selection) bool {
	topNodeScore := extr.getNodeGravityScore(node)
	currentNodeScore := extr.getNodeGravityScore(e)
	threasholdScore := float64(topNodeScore) * 0.08
	if (float64(currentNodeScore) < threasholdScore) && e.Get(0).DataAtom.String() != "td" {
		return false
	}
	return true
}

// we could have long articles that have tons of paragraphs so if we tried to calculate the base score against
// the total text score of those paragraphs it would be unfair. So we need to normalize the score based on the average scoring
// of the paragraphs within the top node. For example if our total score of 10 paragraphs was 1000 but each had an average value of
// 100 then 100 should be our base.
func (extr *ContentExtractor) getSiblingsScore(topNode *goquery.Selection) int {
	base := 100000
	paragraphNumber := 0
	paragraphScore := 0
	nodesToCheck := topNode.Find("p")
	nodesToCheck.Each(func(i int, s *goquery.Selection) {
		textNode := s.Text()
		ws := extr.config.StopWords.StopWordsCount(extr.config.TargetLanguage, textNode)
		highLinkDensity := extr.isHighLinkDensity(s)
		if ws > 2 && !highLinkDensity {
			paragraphNumber++
			paragraphScore += ws
		}
	})
	if paragraphNumber > 0 {
		base = paragraphScore / paragraphNumber
	}
	return base
}

func (extr *ContentExtractor) getSiblingsContent(currentSibling *goquery.Selection, baselinescoreSiblingsPara float64) []*goquery.Selection {
	var ps []*goquery.Selection
	if currentSibling.Get(0).DataAtom.String() == "p" && len(currentSibling.Text()) > 0 {
		ps = append(ps, currentSibling)
		return ps
	}

	potentialParagraphs := currentSibling.Find("p")
	potentialParagraphs.Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if len(text) > 0 {
			ws := extr.config.StopWords.StopWordsCount(extr.config.TargetLanguage, text)
			paragraphScore := ws
			siblingBaselineScore := 0.30
			highLinkDensity := extr.isHighLinkDensity(s)
			score := siblingBaselineScore * baselinescoreSiblingsPara
			if score < float64(paragraphScore) && !highLinkDensity {
				node := new(html.Node)
				node.Type = html.TextNode
				node.Data = text
				node.DataAtom = atom.P
				nodes := make([]*html.Node, 1)
				nodes[0] = node
				newSelection := new(goquery.Selection)
				newSelection.Nodes = nodes
				ps = append(ps, newSelection)
			}
		}

	})
	return ps
}

func (extr *ContentExtractor) walkSiblings(node *goquery.Selection) []*goquery.Selection {
	currentSibling := node.Prev()
	var b []*goquery.Selection
	for currentSibling.Length() != 0 {
		b = append(b, currentSibling)
		previousSibling := currentSibling.Prev()
		currentSibling = previousSibling
	}
	return b
}

// adds any siblings that may have a decent score to this node
func (extr *ContentExtractor) addSiblings(topNode *goquery.Selection) *goquery.Selection {
	if extr.config.Debug {
		log.Println("Starting to add siblings")
	}
	baselinescoreSiblingsPara := extr.getSiblingsScore(topNode)
	results := extr.walkSiblings(topNode)
	for _, currentNode := range results {
		ps := extr.getSiblingsContent(currentNode, float64(baselinescoreSiblingsPara))
		for _, p := range ps {
			nodes := make([]*html.Node, len(topNode.Nodes)+1)
			nodes[0] = p.Get(0)
			for i, node := range topNode.Nodes {
				nodes[i+1] = node
			}
			topNode.Nodes = nodes
		}
	}
	return topNode
}

// PostCleanup removes any divs that looks like non-content, clusters of links, or paras with no gusto
func (extr *ContentExtractor) PostCleanup(targetNode *goquery.Selection) *goquery.Selection {
	if extr.config.Debug {
		log.Println("Starting cleanup Node")
	}
	node := extr.addSiblings(targetNode)
	children := node.Children()
	children.Each(func(i int, s *goquery.Selection) {
		tag := s.Get(0).DataAtom.String()
		if tag != "p" {
			if extr.config.Debug {
				log.Printf("CLEANUP  NODE: %s class: %s\n", extr.config.Parser.Name("id", s), extr.config.Parser.Name("class", s))
			}
			//if extr.isHighLinkDensity(s) || extr.isTableAndNoParaExist(s) || !extr.isNodescoreThresholdMet(node, s) {
			if extr.isHighLinkDensity(s) {
				extr.config.Parser.RemoveNode(s)
				return
			}

			subParagraph := s.Find("p")
			subParagraph.Each(func(j int, e *goquery.Selection) {
				if len(e.Text()) < 25 {
					extr.config.Parser.RemoveNode(e)
				}
			})

			subParagraph2 := s.Find("p")
			if subParagraph2.Length() == 0 && tag != "td" {
				if extr.config.Debug {
					log.Println("Removing node because it doesn't have any paragraphs")
				}
				extr.config.Parser.RemoveNode(s)
			} else {
				if extr.config.Debug {
					log.Println("Not removing TD node")
				}
			}
			return
		}
	})
	return node
}
