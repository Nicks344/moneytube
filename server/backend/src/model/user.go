package model

import (
	"context"
	"time"

	"github.com/meandrewdev/logger"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

func initUsers() {
	usrs, err := GetUsers()
	if err != nil {
		logger.Error(err)
		panic(err)
	}

	for _, user := range usrs {
		if user.IsActivated && user.IsActive {
			ConnectUser(user.Key)
		}
	}
}

func ConnectUser(key string) {
	if _, ok := userDatabases[key]; !ok {
		userDatabases[key] = client.Database("monuytube-user-" + key)
		initIncrements(key)
	}
}

func DisconnectUser(key string) {
	delete(userDatabases, key)
}

type User struct {
	Key            string    `json:"key" bson:"_id"`
	Name           string    `json:"name"`
	CreatedAt      time.Time `json:"createdAt"`
	ActivatedAt    time.Time `json:"activatedAt"`
	IsActivated    bool      `json:"isActivated"`
	IsActive       bool      `json:"isActive"`
	Days           int       `json:"days"`
	DaysLeft       int       `json:"daysLeft" bson:"-"`
	HWID           string    `json:"hwid"`
	EnigmaKey      string    `json:"enigmaKey"`
	DaysReactivate int       `json:"daysReactivate"`
	Version        string    `json:"version"`
}

func (this *User) Init() {
	if !this.IsActivated || !this.IsActive {
		this.DaysLeft = -1
		return
	}

	this.DaysLeft = int((this.ActivatedAt.Add(time.Hour*24*time.Duration(this.Days)).Unix() - time.Now().Unix()) / (24 * 60 * 60))
	if this.DaysLeft < 0 {
		this.IsActive = false
		this.Disactivate()
	}
}

func (this *User) Disactivate() (err error) {
	filter := bson.M{"_id": this.Key}
	ctx, _ := context.WithTimeout(dbCtx, 5*time.Second)
	_, err = getUsersCollection().UpdateOne(ctx, filter, bson.M{"$set": bson.M{"isactive": false}})
	return
}

func GetUser(key string) (user User, err error) {
	filter := bson.M{"_id": key}
	ctx, _ := context.WithTimeout(dbCtx, 5*time.Second)
	err = getUsersCollection().FindOne(ctx, filter).Decode(&user)
	user.Init()
	return
}

func GetUsers() (users []User, err error) {
	users = []User{}
	ctx, _ := context.WithTimeout(dbCtx, 30*time.Second)
	var cur *mongo.Cursor
	cur, err = getUsersCollection().Find(ctx, bson.D{})
	if err != nil {
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result User
		err = cur.Decode(&result)
		if err != nil {
			return
		}
		result.Init()
		users = append(users, result)
	}
	err = cur.Err()
	return
}

func SaveUser(user User) error {
	filter := bson.M{"_id": user.Key}
	ctx, _ := context.WithTimeout(dbCtx, 5*time.Second)
	count, err := getUsersCollection().CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = getUsersCollection().InsertOne(ctx, user)
	} else {
		_, err = getUsersCollection().UpdateOne(ctx, filter, bson.M{"$set": user})
	}

	return err
}

func DeleteUser(key string) (err error) {
	filter := bson.M{"_id": key}
	ctx, _ := context.WithTimeout(dbCtx, 5*time.Second)
	_, err = getUsersCollection().DeleteOne(ctx, filter)
	return
}

func getUsersCollection() *mongo.Collection {
	return database.Collection("users")
}
