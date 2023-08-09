package main

import (
	"flag"
	"log"
	"os"

	"github.com/meandrewdev/logger"
	"github.com/Nicks344/moneytube/server/backend/src/config"
	"github.com/Nicks344/moneytube/server/backend/src/model"
	"github.com/Nicks344/moneytube/server/backend/src/server"
)

func main() {
	isDebug := flag.Bool("debug", false, "is debug")
	flag.Parse()

	logger.Init("logs", "", "")

	if err := config.Init("config", "deploy/docker"); err != nil {
		logger.Error(err)
		log.Fatal(err)
	}

	model.Init()
	initDataPaths()
	server.Serve(*isDebug)
}

func initDataPaths() {
	os.Mkdir("data/reports", 0666)
	os.Mkdir("data/temp", 0666)
	os.Mkdir("data/updates", 0666)
}
