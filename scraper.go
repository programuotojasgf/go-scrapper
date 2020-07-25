package main

import (
	"fmt"
	"github.com/Kamva/mgm/v3"
	"github.com/gocolly/colly/v2"
	"github.com/x/y/data"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"strconv"
)

func main() {
	scrapeReviewsToDatabase()
	//printReviews()
}

func scrapeReviewsToDatabase() {
	reviewCollector := colly.NewCollector()

	reviewCollector.OnHTML(".review-listing", func(e *colly.HTMLElement) {
		externalId, _ := strconv.Atoi(e.ChildAttr("div", "data-review-id"))
		content := e.ChildText(".truncate-content-copy")
		review := data.NewReview(externalId, content)
		mgm.Coll(review).Create(review)
	})

	reviewCollector.OnHTML("a[href][rel='next']", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
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
