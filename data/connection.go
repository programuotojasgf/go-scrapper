package data

import (
	"github.com/Kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"shopify_review_scrapper/config"
)

func init() {
	_ = mgm.SetDefaultConfig(nil, config.Config.DatabaseName, options.Client().ApplyURI(config.Config.ConnectionString))
}
