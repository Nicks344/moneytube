package model

import (
	"context"
	"time"

	"github.com/Nicks344/moneytube/moneytubemodel"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetMacros(key string, name string) (macros moneytubemodel.Macros, err error) {
	filter := bson.M{"_id": name}
	ctx, _ := context.WithTimeout(dbCtx, 5*time.Second)
	err = getMacrosCollection(key).FindOne(ctx, filter).Decode(&macros)
	return
}

func GetMacroses(key string) (macroses []moneytubemodel.Macros, err error) {
	macroses = []moneytubemodel.Macros{}
	ctx, _ := context.WithTimeout(dbCtx, 30*time.Second)
	var cur *mongo.Cursor
	cur, err = getMacrosCollection(key).Find(ctx, bson.D{})
	if err != nil {
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result moneytubemodel.Macros
		err = cur.Decode(&result)
		if err != nil {
			return
		}
		macroses = append(macroses, result)
	}
	err = cur.Err()
	return
}

func SaveMacros(key string, macros moneytubemodel.Macros) error {
	filter := bson.M{"_id": macros.Name}
	ctx, _ := context.WithTimeout(dbCtx, 5*time.Second)
	count, err := getMacrosCollection(key).CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = getMacrosCollection(key).InsertOne(ctx, macros)
	} else {
		_, err = getMacrosCollection(key).UpdateOne(ctx, filter, bson.M{"$set": macros})
	}

	return err
}

func DeleteMacros(key string, name string) (err error) {
	filter := bson.M{"_id": name}
	ctx, _ := context.WithTimeout(dbCtx, 5*time.Second)
	_, err = getMacrosCollection(key).DeleteOne(ctx, filter)
	return
}

func getMacrosCollection(key string) *mongo.Collection {
	return userDatabases[key].Collection("macroses")
}
