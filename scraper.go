package main

import (
	"fmt"
	"github.com/Kamva/mgm/v3"
	"github.com/gocolly/colly/v2"
	"github.com/x/y/data/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		review := models.NewReview(externalId)
		cursor, _ := mgm.Coll(review).Find(mgm.Ctx(), bson.M{"externalId" : externalId}, options.Find().SetLimit(1))
		doesExist := cursor.TryNext(mgm.Ctx())
		if !doesExist {
			mgm.Coll(review).Create(review)
			content := e.ChildText(".truncate-content-copy")
			recordThreeWordPhraseFrequency(content)
		}
	})

	reviewCollector.OnHTML("a[href][rel='next']", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	reviewsUrl := "https://apps.shopify.com/omnisend/reviews"
	reviewCollector.Visit(reviewsUrl)
}

func recordThreeWordPhraseFrequency(content string) {

}

func printReviews() {
	resultReviews := []models.Review{}
	err := mgm.Coll(&models.Review{}).SimpleFind(&resultReviews, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	for _, review := range resultReviews {
		fmt.Printf("Id: %v\n", review.ExternalID)
		fmt.Println()
	}
}
