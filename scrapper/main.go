package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Failed to scrape %s due to %v\n", r.Request.URL.String(), err)
	})

	c.Visit("https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Evolution_of_HTTP")
}
