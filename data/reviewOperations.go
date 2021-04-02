package data

import (
	"github.com/Kamva/mgm/v3"
)

func GetReviewCollection() *mgm.Collection {
	review := &Review{}
	collection := mgm.Coll(review)
	return collection
}