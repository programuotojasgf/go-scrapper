package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func UpsertThreeWordPhraseFrequency(phraseFrequency map[string]int) {
	for phrase, frequency := range phraseFrequency{
		collection := getReviewPhraseCollection()
		options := mgm.UpsertTrueOption()
		filter, updatePipeline := createUpsertFilterAndPipeline(phrase, frequency)
		_, err := collection.UpdateOne(context.Background(), filter, updatePipeline, options)

		if err != nil {
			log.Panic(err)
		}
	}
}

func createUpsertFilterAndPipeline(phrase string, frequency int) (bson.M, []bson.M) {
	filter := bson.M{"phrase": phrase}
	// if it's a new insert, we have to default to 0, because if we read it and it's not there, null + anything is still null
	updateString := fmt.Sprintf(`{ "$set": { "frequency": { "$add": [{ "$ifNull": ["$frequency", 0] }, %d] } } }`, frequency)
	var update bson.M
	json.Unmarshal([]byte(updateString), &update)
	// wrapping the update into a slice makes this a pipeline. You need it to be a pipeline in order
	// for aggregate functions($add and $ifNull) to work. Otherwise they will be treated like fields.
	updatePipeline := []bson.M{update}
	return filter, updatePipeline
}

func getReviewPhraseCollection() *mgm.Collection {
	reviewPhrase := &ReviewPhrase{}
	collection := mgm.Coll(reviewPhrase)
	return collection
}
