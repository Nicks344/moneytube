package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
)

type Report struct {
	ID          primitive.ObjectID `bson:"_id"`
	UserName    string
	Error       string
	Description string
	Resolved    bool
	CreatedAt   time.Time
}

func GetReports() (reports []Report, err error) {
	reports = []Report{}
	ctx, _ := context.WithTimeout(dbCtx, 30*time.Second)
	var cur *mongo.Cursor
	cur, err = getReportsCollection().Find(ctx, bson.D{})
	if err != nil {
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result Report
		err = cur.Decode(&result)
		if err != nil {
			return
		}
		reports = append(reports, result)
	}
	err = cur.Err()
	return
}

func SaveReport(report Report) (string, error) {
	report.ID = primitive.NewObjectID()
	report.CreatedAt = time.Now()

	ctx, _ := context.WithTimeout(dbCtx, 5*time.Second)
	res, err := getReportsCollection().InsertOne(ctx, report)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func DeleteReport(id string) (err error) {
	oID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": oID}
	ctx, _ := context.WithTimeout(dbCtx, 5*time.Second)
	_, err = getReportsCollection().DeleteOne(ctx, filter)
	return
}

func ResolveReport(id string) (err error) {
	oID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": oID}
	ctx, _ := context.WithTimeout(dbCtx, 5*time.Second)
	_, err = getReportsCollection().UpdateOne(ctx, filter, bson.M{"$set": bson.M{"resolved": true}})
	return
}

func getReportsCollection() *mongo.Collection {
	return database.Collection("reports")
}
