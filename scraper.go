package main

import (
	"github.com/Kamva/mgm/v3"
	"github.com/gocolly/colly/v2"
	"github.com/x/y/data"
	"github.com/x/y/phraseCounter"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
)

func main() {
	scrapeReviewsToDatabase()
}

func scrapeReviewsToDatabase() {
	reviewCollector := colly.NewCollector()

	reviewCollector.OnHTML(".review-listing", func(e *colly.HTMLElement) {
		externalId, _ := strconv.Atoi(e.ChildAttr("div", "data-review-id"))
		resultReviews := []data.Review{}
		collection := mgm.Coll(&data.Review{})
		collection.SimpleFind(&resultReviews, bson.M{"externalId" : externalId})
		if len(resultReviews) == 0 {
			collection.Create(data.NewReview(externalId))
			content := e.ChildText(".truncate-content-copy")
			upsertThreeWordPhraseFrequency(phraseCounter.CountThreeWordPhraseFrequency(content))
		}
	})

	reviewCollector.OnHTML("a[href][rel='next']", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	reviewsUrl := "https://apps.shopify.com/omnisend/reviews"
	reviewCollector.Visit(reviewsUrl)
}

func upsertThreeWordPhraseFrequency(phraseFrequency map[string]int) {
	for phrase, frequency := range phraseFrequency{
		reviewPhrase := &data.ReviewPhrase{}
		resultReviewPhrases := []data.ReviewPhrase{}
		collection := mgm.Coll(reviewPhrase)
		collection.SimpleFind(&resultReviewPhrases, bson.M{"phrase" : phrase})
		if (len(resultReviewPhrases) > 0) {
			reviewPhrase = &resultReviewPhrases[0]
			reviewPhrase.Frequency += frequency
			collection.Update(reviewPhrase)
		} else {
			collection.Create(data.NewReviewPhrase(phrase, frequency))
		}
	}
}
