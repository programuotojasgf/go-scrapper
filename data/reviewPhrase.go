package data

import (
	"github.com/Kamva/mgm/v3"
)

// ReviewPhrase is the model for the database ORM
type ReviewPhrase struct {
	mgm.DefaultModel `bson:",inline"`
	Phrase           string `json:"phrase" bson:"phrase"`
	Frequency        int    `json:"frequency" bson:"frequency"`
}

// NewReviewPhrase is a convenience method for creating a new ReviewPhrase
func NewReviewPhrase(phrase string, frequency int) *ReviewPhrase {
	return &ReviewPhrase{
		Phrase:    phrase,
		Frequency: frequency,
	}
}
