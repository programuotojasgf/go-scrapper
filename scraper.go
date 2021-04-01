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

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func canConsumePage(page string, visitedPages *[]string, mutex *sync.Mutex) bool {
	mutex.Lock()
	_, found := Find(*visitedPages, page)
	if !found {
		*visitedPages = append(*visitedPages, page)
	}
	mutex.Unlock()
	return !found
}

func scrapeReviewsToDatabase() {
	reviewCollector := colly.NewCollector(
		colly.Async(true),
		)
	reviewCollector.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 5})

	reviewContentChannel := make(chan string)
	waitGroup := &sync.WaitGroup{}

	reviewCollector.OnHTML(".review-listing", func(e *colly.HTMLElement) {
		externalId, _ := strconv.Atoi(e.ChildAttr("div", "data-review-id"))
		resultReviews := []data.Review{}
		collection := mgm.Coll(&data.Review{})
		collection.SimpleFind(&resultReviews, bson.M{"externalId" : externalId})
		// TODO: In theory, if a review is deleted and a shift happens, this might cause concurrency issues if a review appears in 2 pages at once
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

	visitedPages := make([]string, 0)
	var visitedPagesMutex = &sync.Mutex{}

	reviewCollector.OnHTML(".search-pagination__link", func(e *colly.HTMLElement) {
		nextUrl := e.Attr("href")
		if canConsumePage(nextUrl, &visitedPages, visitedPagesMutex) {
			log.Println("Visiting next page %s", nextUrl)
			e.Request.Visit(nextUrl)
		}
	})

	reviewsUrl := "https://apps.shopify.com/omnisend/reviews?page=1"
	reviewCollector.Visit(reviewsUrl)
	reviewCollector.Wait()

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
