package model

import (
	"context"
	"time"

	"github.com/Nicks344/moneytube/moneytubemodel"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUploadTasks(key string) (tasks []moneytubemodel.UploadTask, err error) {
	return getUploadTasks(key, bson.D{})
}

func GetUploadTasksByDetailsID(key string, detailsID int) (tasks []moneytubemodel.UploadTask, err error) {
	return getUploadTasks(key, bson.M{"detailsid": detailsID})
}

func StopAllTasks(key string) (err error) {
	filter := bson.M{"status": bson.M{"$in": []int{moneytubemodel.UTSStopping, moneytubemodel.UTSInProcess, moneytubemodel.UTSWaiting}}}

	ctx, _ := context.WithTimeout(dbCtx, 5*time.Second)

	_, err = getUploadTaskCollection(key).UpdateMany(ctx, filter, bson.M{"$set": bson.M{"status": moneytubemodel.UTSStopped}})
	return
}

func SaveUploadTask(key string, task moneytubemodel.UploadTask) (id int, err error) {
	filter := bson.M{"_id": task.ID}
	ctx, _ := context.WithTimeout(dbCtx, 5*time.Second)
	if task.ID == 0 {
		id, err = getIncrementID(key, "upload_tasks")
		if err != nil {
			return
		}
		task.ID = id
		_, err = getUploadTaskCollection(key).InsertOne(ctx, task)
	} else {
		id = task.ID
		_, err = getUploadTaskCollection(key).UpdateOne(ctx, filter, bson.M{"$set": task})
	}

	return
}

func DeleteUploadTask(key string, id int) (err error) {
	var task moneytubemodel.UploadTask
	filter := bson.M{"_id": id}
	ctx, _ := context.WithTimeout(dbCtx, 5*time.Second)
	err = getUploadTaskCollection(key).FindOne(ctx, filter).Decode(&task)
	if err != nil {
		return
	}

	ctx, _ = context.WithTimeout(dbCtx, 5*time.Second)
	_, err = getUploadTaskCollection(key).DeleteOne(ctx, filter)
	if err != nil {
		return
	}

	filter = bson.M{"detailsid": task.DetailsID}
	ctx, _ = context.WithTimeout(dbCtx, 5*time.Second)
	count, err := getUploadTaskCollection(key).CountDocuments(ctx, filter)
	if err != nil {
		return err
	}

	if count == 0 {
		err = DeleteUploadData(key, task.DetailsID)
	}

	return
}

func DeleteAllUploadTask(key string) (err error) {
	filter := bson.M{}
	ctx, _ := context.WithTimeout(dbCtx, 5*time.Second)
	_, err = getUploadTaskCollection(key).DeleteMany(ctx, filter)
	if err != nil {
		return
	}
	return DeleteAllUploadData(key)
}

func getUploadTaskCollection(key string) *mongo.Collection {
	return userDatabases[key].Collection("upload_tasks")
}

func getUploadTasks(key string, filter interface{}) (tasks []moneytubemodel.UploadTask, err error) {
	tasks = []moneytubemodel.UploadTask{}
	ctx, _ := context.WithTimeout(dbCtx, 30*time.Second)
	var cur *mongo.Cursor
	cur, err = getUploadTaskCollection(key).Find(ctx, filter)
	if err != nil {
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result moneytubemodel.UploadTask
		err = cur.Decode(&result)
		if err != nil {
			return
		}
		tasks = append(tasks, result)
	}
	err = cur.Err()
	return
}
