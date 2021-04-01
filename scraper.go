package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kamva/mgm/v3"
	"github.com/gocolly/colly/v2"
	"github.com/x/y/data"
	"github.com/x/y/phraseCounter"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"strconv"
	"sync"
)

func main() {
	scrapeReviewsToDatabase()
}

func scrapeReviewsToDatabase() {
	reviewCollector := colly.NewCollector()
	reviewContentChannel := make(chan string)
	waitGroup := &sync.WaitGroup{}

	reviewCollector.OnHTML(".review-listing", func(e *colly.HTMLElement) {
		externalId, _ := strconv.Atoi(e.ChildAttr("div", "data-review-id"))
		resultReviews := []data.Review{}
		collection := mgm.Coll(&data.Review{})
		collection.SimpleFind(&resultReviews, bson.M{"externalId" : externalId})
		if len(resultReviews) == 0 {
			log.Printf("New review %d ! Parsing.\n", externalId)
			collection.Create(data.NewReview(externalId))
			content := e.ChildText(".truncate-content-copy")

			go func(waitGroup *sync.WaitGroup, reviewContentChannel chan<- string) {
				waitGroup.Add(1)
				log.Println("Added another review content for processing!", externalId)
				reviewContentChannel <- content
			}(waitGroup, reviewContentChannel)

			go func(waitGroup *sync.WaitGroup, reviewContentChannel <-chan string) {
				phraseFrequency := phraseCounter.CountThreeWordPhraseFrequency(<-reviewContentChannel)
				upsertThreeWordPhraseFrequency(phraseFrequency)
				log.Println("Processed another review content!", externalId)
				waitGroup.Done()
			}(waitGroup, reviewContentChannel)
		} else {
			log.Printf("Existing review %d !\n", externalId)
		}
	})

	reviewCollector.OnHTML("a[href][rel='next']", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	reviewsUrl := "https://apps.shopify.com/omnisend/reviews"
	reviewCollector.Visit(reviewsUrl)

	log.Println("Waiting for all review contents to be processed.")
	waitGroup.Wait()
	log.Println("Finished processing all review contents.")
}

func upsertThreeWordPhraseFrequency(phraseFrequency map[string]int) {
	for phrase, frequency := range phraseFrequency{
		reviewPhrase := &data.ReviewPhrase{}

		collection := mgm.Coll(reviewPhrase)
		options := mgm.UpsertTrueOption()
		filter := bson.M{"phrase" : phrase}
		updateString := fmt.Sprintf(`{ "$set": { "frequency": { "$add": [{ "$ifNull": ["$frequency", 0] }, %d] } } }`, frequency)
		var update bson.M
		json.Unmarshal([]byte(updateString), &update)
		updatePipeline := []bson.M{update}
		_, err := collection.UpdateOne(context.Background(), filter, updatePipeline, options)

		if err != nil {
			log.Panic(err)
		}
	}
}
