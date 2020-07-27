package data

import (
	"github.com/Kamva/mgm/v3"
)

type ReviewPhrase struct {
	mgm.DefaultModel `bson:",inline"`
	Phrase           string `json:"phrase" bson:"phrase"`
	Frequency        int    `json:"frequency" bson:"frequency"`
}

func NewReviewPhrase(phrase string, frequency int) *ReviewPhrase {
	return &ReviewPhrase{
		Phrase:    phrase,
		Frequency: frequency,
	}
}
