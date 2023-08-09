package model

import (
	"context"
	"time"

	"github.com/Nicks344/moneytube/moneytubemodel"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAccounts(key string) (accounts []moneytubemodel.Account, err error) {
	return GetAccountsByFilter(key, bson.D{})
}

func GetAccountsByFilter(key string, filter interface{}) (accounts []moneytubemodel.Account, err error) {
	accounts = []moneytubemodel.Account{}
	ctx, _ := context.WithTimeout(dbCtx, 30*time.Second)
	var cur *mongo.Cursor
	cur, err = getAccountCollection(key).Find(ctx, filter)
	if err != nil {
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result moneytubemodel.Account
		err = cur.Decode(&result)
		if err != nil {
			return
		}
		accounts = append(accounts, result)
	}
	err = cur.Err()
	return
}

func SaveAccount(key string, account moneytubemodel.Account) (id int, err error) {
	ctx, _ := context.WithTimeout(dbCtx, 5*time.Second)
	filter := bson.M{"_id": account.ID}

	if account.ID == 0 {
		id, err = getIncrementID(key, "accounts")
		if err != nil {
			return
		}
		account.ID = id
		_, err = getAccountCollection(key).InsertOne(ctx, account)
	} else {
		id = account.ID
		_, err = getAccountCollection(key).UpdateOne(ctx, filter, bson.M{"$set": account})
	}

	return
}

func DeleteGroup(key string, group string) (err error) {
	ctx, _ := context.WithTimeout(dbCtx, 5*time.Second)
	filter := bson.M{"group": group}

	_, err = getAccountCollection(key).UpdateMany(ctx, filter, bson.M{"$set": bson.M{"group": ""}})

	return
}

func DeleteAccount(key string, id int) (err error) {
	filter := bson.M{"_id": id}
	ctx, _ := context.WithTimeout(dbCtx, 5*time.Second)
	_, err = getAccountCollection(key).DeleteOne(ctx, filter)
	return
}

func getAccountCollection(key string) *mongo.Collection {
	return userDatabases[key].Collection("accounts")
}
