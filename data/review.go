package data

import (
	"github.com/Kamva/mgm/v3"
)

type Review struct {
	mgm.DefaultModel `bson:",inline"`
	ExternalID       int    `json:"externalId" bson:"externalId"`
}

func NewReview(externalId int) *Review {
	return &Review{
		ExternalID: externalId,
	}
}
