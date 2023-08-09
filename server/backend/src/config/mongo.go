package config

import (
	"fmt"
	"os"
)

func GetMongoConnectURI() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/admin",
		os.Getenv("MONGO_USERNAME"),
		os.Getenv("MONGO_PASSWORD"),
		os.Getenv("MONGO_HOST"),
		os.Getenv("MONGO_PORT"),
	)
}

func GetMongoDBName() string {
	return os.Getenv("MONGO_INITDB_DATABASE")
}
