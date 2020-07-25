package data

import (
	"github.com/Kamva/mgm/v3"
)

type Review struct {
	mgm.DefaultModel `bson:",inline"`
	ExternalID       int    `json:"externalId" bson:"externalId"`
	Content          string `json:"content" bson:"content"`
}

func NewReview(externalId int, content string) *Review {
	return &Review{
		ExternalID: externalId,
		Content:    content,
	}
}
