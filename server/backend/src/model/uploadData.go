package model

import (
	"context"
	"time"

	"github.com/Nicks344/moneytube/client/core/src/utils"
	"github.com/Nicks344/moneytube/moneytubemodel"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUploadDatas(key string) (datas []moneytubemodel.UploadData, err error) {
	datas = []moneytubemodel.UploadData{}
	ctx, _ := context.WithTimeout(dbCtx, 30*time.Second)
	var cur *mongo.Cursor
	cur, err = getUploadDataCollection(key).Find(ctx, bson.D{})
	if err != nil {
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result moneytubemodel.UploadData
		err = cur.Decode(&result)
		if err != nil {
			return
		}
		datas = append(datas, result)
	}
	err = cur.Err()
	return
}

func SaveUploadData(key string, data moneytubemodel.UploadData) (id int, tasks []moneytubemodel.UploadTask, err error) {
	filter := bson.M{"_id": data.ID}
	ctx, _ := context.WithTimeout(dbCtx, 5*time.Second)
	if data.ID == 0 {
		id, err = getIncrementID(key, "upload_datas")
		if err != nil {
			return
		}
		data.ID = id
		_, err = getUploadDataCollection(key).InsertOne(ctx, data)
	} else {
		id = data.ID
		_, err = getUploadDataCollection(key).UpdateOne(ctx, filter, bson.M{"$set": data})
		if err != nil {
			return
		}
		tasks, err = GetUploadTasksByDetailsID(key, id)
		if err != nil {
			return
		}
		for i := range tasks {
			tasks[i].Status = moneytubemodel.UTSStopped
			tasks[i].ErrorMessage = ""
			tasks[i].Count = utils.RandRange(data.UploadCountFrom, data.UploadCountTo)
			tasks[i].Scheduler = moneytubemodel.UploadTaskScheduler{
				Enabled:         data.Scheduler.Enabled,
				SecondStartTime: data.Scheduler.StartTime,
				Progress:        0,
			}
			if _, err = SaveUploadTask(key, tasks[i]); err != nil {
				return
			}
		}
		return
	}

	tasks = make([]moneytubemodel.UploadTask, len(data.AccountIDs), len(data.AccountIDs))
	for i, accID := range data.AccountIDs {
		task := moneytubemodel.UploadTask{
			AccountID: accID,
			DetailsID: data.ID,
			Status:    moneytubemodel.UTSStopped,
			Count:     utils.RandRange(data.UploadCountFrom, data.UploadCountTo),
			Scheduler: moneytubemodel.UploadTaskScheduler{
				Enabled:         data.Scheduler.Enabled,
				SecondStartTime: data.Scheduler.StartTime,
				Progress:        0,
			},
		}
		var taskID int
		if taskID, err = SaveUploadTask(key, task); err != nil {
			return
		}
		task.ID = taskID
		tasks[i] = task
	}

	return
}

func DeleteUploadData(key string, id int) (err error) {
	filter := bson.M{"_id": id}
	ctx, _ := context.WithTimeout(dbCtx, 5*time.Second)
	_, err = getUploadDataCollection(key).DeleteOne(ctx, filter)
	return
}

func DeleteAllUploadData(key string) (err error) {
	filter := bson.M{}
	ctx, _ := context.WithTimeout(dbCtx, 5*time.Second)
	_, err = getUploadDataCollection(key).DeleteMany(ctx, filter)
	return
}

func getUploadDataCollection(key string) *mongo.Collection {
	return userDatabases[key].Collection("upload_datas")
}
