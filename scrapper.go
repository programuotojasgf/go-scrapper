package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector()

	c.OnHTML(".review-listing", func(e *colly.HTMLElement) {
		id := e.ChildAttr("div", "data-review-id")
		fmt.Printf("id: %v\n", id)

		shortContent := e.ChildText(".truncate-content-copy")
		fmt.Printf("Short content: %v\n", shortContent)

		fmt.Println()
	})

	reviewsUrl := "https://apps.shopify.com/omnisend/reviews"
	c.Visit(reviewsUrl)
}
