package main

import (
	"context"
	"github.com/gocolly/colly/v2"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"shopify_review_scrapper/config"
	"shopify_review_scrapper/data"
	"shopify_review_scrapper/phraseCounter"
	"strconv"
	"sync"
)

func main() {
	reviewCollector := createConfiguredReviewCollector()
	scrapeReviewsToDatabase(reviewCollector, config.Config.ReviewsUrlFirstPage)
}

func createConfiguredReviewCollector() *colly.Collector {
	reviewCollector := colly.NewCollector(
		colly.Async(true),
	)
	reviewCollector.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: config.Config.WebsiteVisitorParallelismLimit})
	return reviewCollector
}

func scrapeReviewsToDatabase(reviewCollector *colly.Collector, reviewsURL string) {
	reviewProcessingWaitGroup := setupReviewCollector(reviewCollector)

	reviewCollector.Visit(reviewsURL)
	reviewCollector.Wait()

	log.Println("Waiting for all review contents to be processed.")
	reviewProcessingWaitGroup.Wait()
	log.Println("Finished processing all review contents.")
}

func setupReviewCollector(reviewCollector *colly.Collector) *sync.WaitGroup {
	reviewProcessingWaitGroup := setupScrappingAndParsingReviews(reviewCollector)
	setupVisitingPages(reviewCollector)
	return reviewProcessingWaitGroup
}

func setupScrappingAndParsingReviews(reviewCollector *colly.Collector) *sync.WaitGroup {
	reviewContentChannel := make(chan string)
	waitGroup := &sync.WaitGroup{}

	reviewCollector.OnHTML(".review-listing", func(e *colly.HTMLElement) {
		externalID, _ := strconv.Atoi(e.ChildAttr("div", "data-review-id"))
		collection := data.GetReviewCollection()
		count, err := collection.CountDocuments(context.Background(), bson.M{"externalID": externalID})

		if err != nil {
			panic(err)
		}

		// TODO: In theory, if a review is deleted and a shift happens, this might cause concurrency issues if the same review appears in 2 pages at once
		if count == 0 {
			log.Printf("New review %d ! Parsing.\n", externalID)
			collection.Create(data.NewReview(externalID))
			content := e.ChildText(".truncate-content-copy")
			go addReviewToReviewContentChannel(externalID, content, waitGroup, reviewContentChannel)
			go consumeReviewFromReviewContentChannel(externalID, waitGroup, reviewContentChannel)
		} else {
			log.Printf("Existing review %d !\n", externalID)
		}
	})
	return waitGroup
}

func consumeReviewFromReviewContentChannel(externalID int, waitGroup *sync.WaitGroup, reviewContentChannel <-chan string) {
	phraseFrequency := phraseCounter.CountThreeWordPhraseFrequency(<-reviewContentChannel)
	data.UpsertThreeWordPhraseFrequency(phraseFrequency)
	log.Println("Processed another review content!", externalID)
	waitGroup.Done()
}

func addReviewToReviewContentChannel(externalID int, content string, waitGroup *sync.WaitGroup, reviewContentChannel chan<- string) {
	waitGroup.Add(1)
	log.Println("Added another review content for processing!", externalID)
	reviewContentChannel <- content
}

// We are using a slice to track already visited pages, to ensure each page is only visited once.
func setupVisitingPages(reviewCollector *colly.Collector) {
	visitedPages := make([]string, 0)
	var visitedPagesMutex = &sync.Mutex{}

	reviewCollector.OnHTML(".search-pagination__link", func(e *colly.HTMLElement) {
		pageURL := e.Attr("href")
		if canConsumePage(pageURL, &visitedPages, visitedPagesMutex) {
			log.Println("Visiting next page ", pageURL)
			e.Request.Visit(pageURL)
		}
	})
}
