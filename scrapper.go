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
	testDatabase()
	return

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

func testDatabase() {
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

//func testDatabase() {
//	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://test_user:4TKpw5P9tPCjscnO@cluster0.ufryp.mongodb.net/scrapper?retryWrites=true&w=majority"))
//	if err != nil {
//		log.Fatal(err)
//	}
//	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
//	err = client.Connect(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer client.Disconnect(ctx)
//
//	scrapperDatabase := client.Database("scrapper")
//	reviewsCollection := scrapperDatabase.Collection("reviews")
//
//	cursor, findError := reviewsCollection.Find(
//		context.Background(),
//		bson.D{{}},
//		options.Find(),
//	)
//
//	if findError != nil {
//		log.Fatal(findError)
//	}
//
//	for cursor.Next(ctx) {
//		var review Review
//		err := cursor.Decode(&review)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		fmt.Println(review)
//	}
//}
