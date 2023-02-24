package main

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"github.com/gocolly/colly/v2/extensions"
)

func main() {
	// Instantiate default collector
	c := colly.NewCollector(

		colly.Debugger(&debug.LogDebugger{}),

		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("https://www.wikipedia.org/"),
		colly.MaxDepth(2),
		colly.AllowURLRevisit(),
		colly.URLFilters(
			regexp.MustCompile("https://www.wikipedia.org/"),
		),
		colly.UserAgent("Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1)"),
	)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*wikipedia.org.*",
		Parallelism: 1,
		RandomDelay: 5 * time.Second,
	})

	extensions.Referer(c)
	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		fmt.Printf("before link visit %s", link)

		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		c.Visit(e.Request.AbsoluteURL(link))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("response ")
		log.Printf("%s\n", bytes.Replace(r.Body, []byte("\n"), nil, -1))
	})

	// Start scraping on https://hackerspaces.org
	err := c.Visit("https://www.wikipedia.org/")
	if err != nil {
		fmt.Println(err)
	}
}
