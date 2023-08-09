package model

import (
	"context"
	"time"

	"github.com/Nicks344/moneytube/moneytubemodel"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUploadDataTemplates(key string) (templates []moneytubemodel.UploadDataTemplate, err error) {
	templates = []moneytubemodel.UploadDataTemplate{}
	ctx, _ := context.WithTimeout(dbCtx, 30*time.Second)
	var cur *mongo.Cursor
	cur, err = getUploadDataTemplateCollection(key).Find(ctx, bson.D{})
	if err != nil {
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result moneytubemodel.UploadDataTemplate
		err = cur.Decode(&result)
		if err != nil {
			return
		}
		templates = append(templates, result)
	}
	err = cur.Err()
	return
}

func SaveUploadDataTemplate(key string, data moneytubemodel.UploadDataTemplate) (err error) {
	ctx, _ := context.WithTimeout(dbCtx, 5*time.Second)
	filter := bson.M{"_id": data.Label}

	var count int64
	count, err = getUploadDataTemplateCollection(key).CountDocuments(ctx, filter)
	if err != nil {
		return
	}

	if count == 0 {
		_, err = getUploadDataTemplateCollection(key).InsertOne(ctx, data)
	} else {
		_, err = getUploadDataTemplateCollection(key).UpdateOne(ctx, filter, bson.M{"$set": data})
	}
	return
}

func DeleteUploadDataTemplate(key string, id string) (err error) {
	filter := bson.M{"_id": id}
	ctx, _ := context.WithTimeout(dbCtx, 5*time.Second)
	_, err = getUploadDataTemplateCollection(key).DeleteOne(ctx, filter)
	return
}

func DeleteAllUploadDataTemplate(key string) (err error) {
	filter := bson.M{}
	ctx, _ := context.WithTimeout(dbCtx, 5*time.Second)
	_, err = getUploadDataTemplateCollection(key).DeleteMany(ctx, filter)
	return
}

func getUploadDataTemplateCollection(key string) *mongo.Collection {
	return userDatabases[key].Collection("upload_data_templates")
}
