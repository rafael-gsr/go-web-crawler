// Package crawler contains the crawler implementation, thinking on a modular feature
package crawler

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	datastructure "github.com/rafael-gsr/go-web-crawler/internal/data_structure"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Crawler struct {
	links   datastructure.Set[string]
	visited datastructure.Set[string]
	emails  datastructure.Set[string]
	skipped datastructure.Set[string]
	current string

	shouldCrawEmails     bool
	shouldKeepSameDomain bool
}

func (c *Crawler) Emails() *datastructure.Set[string]  { return &c.emails }
func (c *Crawler) Links() *datastructure.Set[string]   { return &c.links }
func (c *Crawler) Visited() *datastructure.Set[string] { return &c.visited }
func (c *Crawler) Current() string                     { return c.current }

func (c *Crawler) SetEmailCrawling(on bool) {
	c.shouldCrawEmails = on
}

func (c *Crawler) SetKeepSameDomain(on bool) {
	c.shouldKeepSameDomain = on
}

func (c *Crawler) findLinks(node *html.Node) (err error) {
	isLinkNode := node.DataAtom == atom.A

	if isLinkNode {
		for _, attrbt := range node.Attr {

			isLinkAttribute := attrbt.Key == "href"

			if isLinkAttribute {
				link := attrbt.Val
				isExternalLink := len(link) > 3 && link[:4] == "http"

				if !isExternalLink {
					completeInternalPath, err := url.JoinPath(c.current, link)
					if err != nil {
						continue
					}
					link = completeInternalPath
				}

				isUnmapped := !(c.visited.Contains(link) || c.links.Contains(link) || c.skipped.Contains(link))
				isDomainAcceptable := !c.shouldKeepSameDomain || (c.shouldKeepSameDomain && strings.Contains(link, c.current))

				if isUnmapped && isDomainAcceptable {
					c.links.Add(link)
				}

			}
		}
	}
	return
}

func (c *Crawler) findEmails(node *html.Node) (err error) {
	emailRegex, err := regexp.Compile("^[\\w\\-\\.]+@([\\w-]+\\.)+[\\w-]{2,}$")
	if err != nil {
		return
	}

	nodeContent := node.Data
	matches := emailRegex.FindAllString(nodeContent, -1)

	for _, email := range matches {
		c.emails.Add(email)
	}

	return
}

func (c Crawler) printList(list datastructure.Set[string]) {
	for _, item := range list.GetList() {
		fmt.Println(item)
	}
}

func (c *Crawler) PrintSkipped() {
	c.printList(c.skipped)
}

func (c *Crawler) PrintVisited() {
	c.printList(c.visited)
}

func (c *Crawler) PrintEmails() {
	c.printList(c.emails)
}

func (c *Crawler) PrintLinks() {
	c.printList(c.links)
}

func (c *Crawler) Crawl() {
	res, err := http.Get(c.current)
	if err != nil {
		return
	}

	parsedHTML, err := html.Parse(res.Body)
	if err != nil {
		return
	}

	for des := range parsedHTML.Descendants() {

		err := c.findLinks(des)
		if err != nil {
			return
		}

		if c.shouldCrawEmails {
			err = c.findEmails(des)
			if err != nil {
				return
			}
		}

	}

	c.visited.Add(c.current)
}

func (c *Crawler) SetupNext() {
	next := c.links.GetFirst()
	c.current = next
	c.links.Remove(next)
}

func (c *Crawler) Skip() {
	skipped := c.current
	c.skipped.Add(skipped)
	c.SetupNext()
}

func NewCrawler(initialPage string, shouldKeepSameDomain bool, shouldFindEmails bool) *Crawler {
	links := datastructure.NewSet[string]()
	links.Add(initialPage)

	crawler := Crawler{
		links:                links,
		visited:              datastructure.NewSet[string](),
		emails:               datastructure.NewSet[string](),
		skipped:              datastructure.NewSet[string](),
		shouldCrawEmails:     shouldFindEmails,
		shouldKeepSameDomain: shouldKeepSameDomain,
	}

	return &crawler
}

func GetPageLinks(pageUrl string) (links []string, err error) {
	return
}
