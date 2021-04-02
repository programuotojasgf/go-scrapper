package data

import (
	"github.com/Kamva/mgm/v3"
)

// Review is the model for the database ORM
type Review struct {
	mgm.DefaultModel `bson:",inline"`
	ExternalID       int `json:"externalId" bson:"externalId"`
}

// NewReview is a convenience method for creating a new Review
func NewReview(externalID int) *Review {
	return &Review{
		ExternalID: externalID,
	}
}
