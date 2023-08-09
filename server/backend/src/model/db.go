package model

import (
	"context"

	"github.com/meandrewdev/logger"
	"github.com/Nicks344/moneytube/server/backend/src/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var database *mongo.Database
var userDatabases = map[string]*mongo.Database{}
var dbCtx context.Context

func Init() {
	dbCtx = context.Background()
	var err error
	client, err = mongo.Connect(dbCtx, options.Client().ApplyURI(config.GetMongoConnectURI()))
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	database = client.Database(config.GetMongoDBName())

	initUsers()
}
