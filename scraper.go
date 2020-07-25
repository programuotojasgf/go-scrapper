package main

import (
	"fmt"
	"github.com/Kamva/mgm/v3"
	"github.com/gocolly/colly/v2"
	"github.com/x/y/data"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func main() {
	printReviews()
	return

	ScrapeReviewsToDatabase()
}

func ScrapeReviewsToDatabase() {
	reviewCollector := colly.NewCollector()

	reviewCollector.OnHTML(".review-listing", func(e *colly.HTMLElement) {
		id := e.ChildAttr("div", "data-review-id")
		fmt.Printf("id: %v\n", id)

		shortContent := e.ChildText(".truncate-content-copy")
		fmt.Printf("Short content: %v\n", shortContent)

		fmt.Println()
	})

	reviewsUrl := "https://apps.shopify.com/omnisend/reviews"
	reviewCollector.Visit(reviewsUrl)
}

func printReviews() {
	resultReviews := []data.Review{}
	err := mgm.Coll(&data.Review{}).SimpleFind(&resultReviews, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	for _, review := range resultReviews {
		fmt.Printf("Id: %v\n", review.ExternalID)
		fmt.Printf("Content: %v\n", review.Content)
		fmt.Println()
	}
}
