package goose

import (
	"container/list"
	"github.com/advancedlogic/goquery"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"log"
	"regexp"
	"strings"
)

// Cleaner removes menus, ads, sidebars, etc. and leaves the main content
type Cleaner struct {
	config Configuration
}

// NewCleaner returns a new instance of a Cleaner
func NewCleaner(config Configuration) Cleaner {
	return Cleaner{
		config: config,
	}
}

var divToPElementsPattern = regexp.MustCompile("<(a|blockquote|dl|div|img|ol|p|pre|table|ul)")
var tabsRegEx, _ = regexp.Compile("\\t|^\\s+$]")
var removeNodesRegEx = regexp.MustCompile("" +
	"PopularQuestions|" +
	"[Cc]omentario|" +
	"[Ff]ooter|" +
	"^fn$|" +
	"^inset$|" +
	"^print$|" +
	"^scroll$|" +
	"^side$|" +
	"^side_|" +
	"^widget$|" +
	"ajoutVideo|" +
	"articleheadings|" +
	"author-dropdown|" +
	"blog-pager|" +
	"breadcrumbs|" +
	"byline|" +
	"cabecalho|" +
	"cnnStryHghLght|" +
	"cnn_html_slideshow|" +
	"cnn_strycaptiontxt|" +
	"cnn_strylftcntnt|" +
	"cnn_stryspcvbx|" +
	"combx|" +
	"comment|" +
	"communitypromo|" +
	"contact|" +
	"contentTools2|" +
	"controls|" +
	"^date$|" +
	"detail_new_|" +
	"detail_related_|" +
	"figcaption|" +
	"footnote|" +
	"foot|" +
	"header|" +
	"img_popup_single|" +
	"js_replies|" +
	"[Kk]ona[Ff]ilter|" +
	"leading|" +
	"legende|" +
	"links|" +
	"mediaarticlerelated|" +
	"menucontainer|" +
	"meta$|" +
	"navbar|" +
	"pagetools|" +
	"popup|" +
	"post-attributes|" +
	"post-title|" +
	"relacionado|" +
	"retweet|" +
	"runaroundLeft|" +
	"shoutbox|" +
	"site_nav|" +
	"socialNetworking|" +
	"social_|" +
	"socialnetworking|" +
	"socialtools|" +
	"sponsor|" +
	"sub_nav|" +
	"subscribe|" +
	"tag_|" +
	"tags|" +
	"the_answers|" +
	"timestamp|" +
	"tools|" +
	"vcard|" +
	"welcome_form|" +
	"wp-caption-text")
var captionsRegEx = regexp.MustCompile("^caption$")
var googleRegEx = regexp.MustCompile(" google ")
var moreRegEx = regexp.MustCompile("^[^entry-]more.*$")
var facebookRegEx = regexp.MustCompile("[^-]facebook")
var facebookBroadcastingRegEx = regexp.MustCompile("facebook-broadcasting")
var twitterRegEx = regexp.MustCompile("[^-]twitter")

func (c *Cleaner) clean(article *Article) *goquery.Document {
	if c.config.debug {
		log.Println("Starting cleaning phase with Cleaner")
	}
	docToClean := article.Doc
	docToClean = c.cleanArticleTags(docToClean)
	docToClean = c.cleanEMTags(docToClean)
	docToClean = c.dropCaps(docToClean)
	docToClean = c.removeScriptsStyle(docToClean)
	docToClean = c.cleanBadTags(docToClean)
	docToClean = c.cleanFooter(docToClean)
	docToClean = c.cleanAside(docToClean)
	docToClean = c.removeNodesRegEx(docToClean, captionsRegEx)
	docToClean = c.removeNodesRegEx(docToClean, googleRegEx)
	docToClean = c.removeNodesRegEx(docToClean, moreRegEx)
	docToClean = c.removeNodesRegEx(docToClean, facebookRegEx)
	docToClean = c.removeNodesRegEx(docToClean, facebookBroadcastingRegEx)
	docToClean = c.removeNodesRegEx(docToClean, twitterRegEx)
	docToClean = c.cleanParaSpans(docToClean)

	docToClean = c.convertDivsToParagraphs(docToClean, "div")
	docToClean = c.convertDivsToParagraphs(docToClean, "span")
	docToClean = c.convertDivsToParagraphs(docToClean, "article")
	docToClean = c.convertDivsToParagraphs(docToClean, "pre")

	return docToClean
}

func (c *Cleaner) cleanArticleTags(doc *goquery.Document) *goquery.Document {
	tags := [3]string{"id", "name", "class"}
	articles := doc.Find("article")
	articles.Each(func(i int, s *goquery.Selection) {
		for _, tag := range tags {
			c.config.parser.delAttr(s, tag)
		}
	})
	return doc
}

func (c *Cleaner) cleanEMTags(doc *goquery.Document) *goquery.Document {
	ems := doc.Find("em")
	ems.Each(func(i int, s *goquery.Selection) {
		images := s.Find("img")
		if images.Length() == 0 {
			c.config.parser.dropTag(s)
		}
	})
	if c.config.debug {
		log.Printf("Cleaning %d EM tags\n", ems.Size())
	}
	return doc
}

func (c *Cleaner) cleanFooter(doc *goquery.Document) *goquery.Document {
	footer := doc.Find("footer")
	footer.Each(func(i int, s *goquery.Selection) {
		c.config.parser.removeNode(s)
	})
	return doc
}

func (c *Cleaner) cleanAside(doc *goquery.Document) *goquery.Document {
	aside := doc.Find("aside")
	aside.Each(func(i int, s *goquery.Selection) {
		c.config.parser.removeNode(s)
	})
	return doc
}

func (c *Cleaner) cleanCites(doc *goquery.Document) *goquery.Document {
	cites := doc.Find("cite")
	cites.Each(func(i int, s *goquery.Selection) {
		c.config.parser.removeNode(s)
	})
	return doc
}

func (c *Cleaner) cleanDivs(doc *goquery.Document) *goquery.Document {
	frames := make(map[string]int)
	framesNodes := make(map[string]*list.List)
	divs := doc.Find("div")
	divs.Each(func(i int, s *goquery.Selection) {
		children := s.Children()
		if children.Size() == 0 {
			text := s.Text()
			text = strings.Trim(text, " ")
			text = strings.Trim(text, "\t")
			text = strings.ToLower(text)
			frames[text]++
			if framesNodes[text] == nil {
				framesNodes[text] = list.New()
			}
			framesNodes[text].PushBack(s)
		}
	})
	for text, freq := range frames {
		if freq > 1 {
			selections := framesNodes[text]
			for s := selections.Front(); s != nil; s = s.Next() {
				selection := s.Value.(*goquery.Selection)
				c.config.parser.removeNode(selection)
			}
		}
	}
	return doc
}

func (c *Cleaner) dropCaps(doc *goquery.Document) *goquery.Document {
	items := doc.Find("span")
	count := 0 //remove
	items.Each(func(i int, s *goquery.Selection) {
		attribute, exists := s.Attr("class")
		if exists && (strings.Contains(attribute, "dropcap") || strings.Contains(attribute, "drop_cap")) {
			count++
			c.config.parser.dropTag(s)
		}
	})
	if c.config.debug {
		log.Printf("Cleaning %d dropcap tags\n", count)
	}
	return doc
}

func (c *Cleaner) removeScriptsStyle(doc *goquery.Document) *goquery.Document {
	if c.config.debug {
		log.Println("Starting to remove script tags")
	}
	scripts := doc.Find("script,noscript,style")
	scripts.Each(func(i int, s *goquery.Selection) {
		c.config.parser.removeNode(s)
	})
	if c.config.debug {
		log.Printf("Removed %d script and style tags\n", scripts.Size())
	}

	//remove comments :) How????
	return doc
}

func (c *Cleaner) matchNodeRegEx(attribute string, pattern *regexp.Regexp) bool {
	return pattern.MatchString(attribute)
}

func (c *Cleaner) removeNodesRegEx(doc *goquery.Document, pattern *regexp.Regexp) *goquery.Document {
	selectors := [3]string{"id", "class", "name"}
	for _, selector := range selectors {
		naughtyList := doc.Find("*[" + selector + "]")
		cont := 0
		naughtyList.Each(func(i int, s *goquery.Selection) {
			attribute, _ := s.Attr(selector)
			if c.matchNodeRegEx(attribute, pattern) {
				cont++
				c.config.parser.removeNode(s)
			}
		})

		if c.config.debug {
			log.Printf("regExRemoveNodes %d %s elements found against pattern %s\n", cont, selector, pattern.String())
		}
	}
	return doc
}

func (c *Cleaner) cleanBadTags(doc *goquery.Document) *goquery.Document {
	body := doc.Find("body")
	children := body.Children()
	selectors := []string{"id", "class", "name"}
	for _, selector := range selectors {
		children.Each(func(i int, s *goquery.Selection) {
			naughtyList := s.Find("*[" + selector + "]")
			cont := 0
			naughtyList.Each(func(j int, e *goquery.Selection) {
				attribute, _ := e.Attr(selector)
				if c.matchNodeRegEx(attribute, removeNodesRegEx) {
					if c.config.debug {

						log.Printf("Cleaning: Removing node with %s: %s\n", selector, c.config.parser.name(selector, e))
					}
					c.config.parser.removeNode(e)
					cont++
				}
			})
			if c.config.debug && cont > 0 {
				log.Printf("%d naughty %s elements found", cont, selector)
			}
		})
	}
	return doc
}

func (c *Cleaner) cleanParaSpans(doc *goquery.Document) *goquery.Document {
	spans := doc.Find("span")
	spans.Each(func(i int, s *goquery.Selection) {
		parent := s.Parent()
		if parent != nil && parent.Length() > 0 && parent.Get(0).DataAtom == atom.P {
			node := s.Get(0)
			node.Data = s.Text()
			node.Type = html.TextNode
		}
	})
	return doc
}

func (c *Cleaner) getFlushedBuffer(fragment string) []*html.Node {
	var output []*html.Node
	reader := strings.NewReader(fragment)
	document, _ := html.Parse(reader)
	body := document.FirstChild.LastChild
	for c := body.FirstChild; c != nil; c = c.NextSibling {
		output = append(output, c)
		c.Parent = nil
		c.PrevSibling = nil
	}

	for _, o := range output {
		o.NextSibling = nil
	}
	return output
}

func (c *Cleaner) replaceWithPara(div *goquery.Selection) {
	if div.Size() > 0 {
		node := div.Get(0)
		node.Data = atom.P.String()
		node.DataAtom = atom.P
	}
}

func (c *Cleaner) tabsAndNewLinesReplacements(text string) string {
	text = strings.Replace(text, "\n", "\n\n", -1)
	text = tabsRegEx.ReplaceAllString(text, "")
	return text
}

func (c *Cleaner) convertDivsToParagraphs(doc *goquery.Document, domType string) *goquery.Document {
	if c.config.debug {
		log.Println("Starting to replace bad divs...")
	}
	badDivs := 0
	convertedTextNodes := 0
	divs := doc.Find(domType)

	divs.Each(func(i int, div *goquery.Selection) {
		divHTML, _ := div.Html()
		if divToPElementsPattern.Match([]byte(divHTML)) {
			c.replaceWithPara(div)
			badDivs++
		} else {
			var replacementText []string
			nodesToRemove := list.New()
			children := div.Contents()
			if c.config.debug {
				log.Printf("Found %d children of div\n", children.Size())
			}
			children.EachWithBreak(func(i int, kid *goquery.Selection) bool {
				text := kid.Text()
				kidNode := kid.Get(0)
				tag := kidNode.Data
				if tag == text {
					tag = "#text"
				}
				if tag == "#text" {
					text = strings.Replace(text, "\n", "", -1)
					text = tabsRegEx.ReplaceAllString(text, "")
					if text == "" {
						return true
					}
					if len(text) > 1 {
						prev := kidNode.PrevSibling
						if c.config.debug {
							log.Printf("PARENT CLASS: %s NODENAME: %s\n", c.config.parser.name("class", div), tag)
							log.Printf("TEXTREPLACE: %s\n", strings.Replace(text, "\n", "", -1))
						}
						if prev != nil && prev.DataAtom == atom.A {
							nodeSelection := kid.HasNodes(prev)
							html, _ := nodeSelection.Html()
							replacementText = append(replacementText, html)
							if c.config.debug {
								log.Printf("SIBLING NODENAME ADDITION: %s TEXT: %s\n", prev.Data, html)
							}
						}
						replacementText = append(replacementText, text)
						nodesToRemove.PushBack(kidNode)
						convertedTextNodes++
					}

				}
				return true
			})

			newNode := new(html.Node)
			newNode.Type = html.ElementNode
			newNode.Data = strings.Join(replacementText, "")
			newNode.DataAtom = atom.P
			div.First().AddNodes(newNode)

			for s := nodesToRemove.Front(); s != nil; s = s.Next() {
				node := s.Value.(*html.Node)
				if node != nil && node.Parent != nil {
					node.Parent.RemoveChild(node)
				}
			}
		}
	})
	if c.config.debug {
		log.Printf("Found %d total divs with %d bad divs replaced and %d textnodes converted inside divs", divs.Size(), badDivs, convertedTextNodes)
	}
	return doc

}
