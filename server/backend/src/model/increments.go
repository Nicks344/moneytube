package model

import (
	"context"
	"github.com/meandrewdev/logger"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var increments = map[string]map[string]*Increment{}

func initIncrements(key string) {
	increments[key] = map[string]*Increment{}

	ctx, _ := context.WithTimeout(dbCtx, 30*time.Second)
	cur, err := getIncrementsCollection(key).Find(ctx, bson.D{})
	if err != nil {
		logger.Error(err)
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var entry Increment
		err := cur.Decode(&entry)
		if err != nil {
			logger.Error(err)
			break
		}
		increments[key][entry.Name] = &entry
	}
}

type Increment struct {
	Name  string `json:"name" bson:"_id"`
	Value int    `json:"value" bson:"value"`
}

func (inc *Increment) Inc(key string) (int, error) {
	ctx, _ := context.WithTimeout(dbCtx, 5*time.Second)
	_, err := getIncrementsCollection(key).UpdateOne(ctx, bson.M{"_id": inc.Name}, bson.M{"$inc": bson.M{"value": 1}})
	if err != nil {
		logger.Error(err)
		return 0, err
	}
	inc.Value++
	return inc.Value, nil
}

func getIncrementID(key string, collectionName string) (id int, err error) {
	ctx, _ := context.WithTimeout(dbCtx, 5*time.Second)
	inc, ok := increments[key][collectionName]
	if !ok {
		inc = &Increment{
			Name:  collectionName,
			Value: 1,
		}
		_, err = getIncrementsCollection(key).InsertOne(ctx, inc)
		if err != nil {
			logger.Error(err)
			return
		}
		increments[key][collectionName] = inc
		return 1, nil
	}
	return inc.Inc(key)
}

func getIncrementsCollection(key string) *mongo.Collection {
	return userDatabases[key].Collection("increments")
}
