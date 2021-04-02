package data

import (
	"github.com/Kamva/mgm/v3"
)

// GetReviewCollection is a convenience method of getting the review collection
func GetReviewCollection() *mgm.Collection {
	review := &Review{}
	collection := mgm.Coll(review)
	return collection
}
