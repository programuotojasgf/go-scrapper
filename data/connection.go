package data

import (
	"github.com/Kamva/mgm/v3"
	"github.com/x/y/config"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	_ = mgm.SetDefaultConfig(nil, "scrapper", options.Client().ApplyURI(config.Config.ConnectionString))
}
