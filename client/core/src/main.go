package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/meandrewdev/logger"
	"github.com/Nicks344/moneytube/client/core/src/clibridge"
	"github.com/Nicks344/moneytube/client/core/src/config"
	"github.com/Nicks344/moneytube/client/core/src/license"
	"github.com/Nicks344/moneytube/client/core/src/model"
	"github.com/Nicks344/moneytube/client/core/src/server"
	"github.com/Nicks344/moneytube/client/core/src/server/gqlserver"
	"github.com/Nicks344/moneytube/client/core/src/serverAPI"
	"github.com/Nicks344/moneytube/client/core/src/uibridge"
	"github.com/Nicks344/moneytube/client/core/src/utils/update"
	"github.com/Nicks344/moneytube/client/core/src/utils/videoeditor"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(fmt.Errorf("PANIC: %v", err))
		}
	}()

	videoeditor.SetPaths(config.GetFFmpegBin(), config.GetFFprobeBin())

	gqlPort := flag.Int("gql", 10000, "gqlPort")
	rpcClientPort := flag.Int("rpcClient", 10010, "rpcClientPort")
	rpcUIPort := flag.Int("rpcUI", 10020, "rpcUI")
	rpcCLIPort := flag.Int("rpcCLI", 10030, "rpcCLI")
	logPath := flag.String("log", "logs", "logs path")
	flag.Parse()

	logger.Init(*logPath, "", "")
	logger.SetStdout(true)
	model.Init()

	gqlserver.Init()
	uibridge.Connect(*rpcUIPort)
	go uibridge.Serve(*rpcClientPort)
	go clibridge.Serve(*rpcCLIPort)
	uibridge.OnActivate(func(key string) error {
		err := license.Register(key)
		if err != nil {
			logger.Error(err)
			return err
		}
		launch(*gqlPort)
		return nil
	})

	launch(*gqlPort)

	for {
		fmt.Scanln()
	}
}

func launch(gqlPort int) {
	if license.Check() {
		if checkUpdate() {
			return
		}

		uibridge.Endpoint.Launch(config.GetApiKey(), config.GetVersion())
		serverAPI.StopAllTasks()
		go server.Start(gqlPort)
	} else {
		uibridge.Endpoint.Login()
	}
}

func checkUpdate() bool {
	version, err := update.Check()
	if err != nil {
		uibridge.Endpoint.OnUpdateError(err.Error())
		return true
	}
	if version != "" {
		uibridge.Endpoint.OnUpdating(version)
		ctx, cancel := context.WithCancel(context.Background())
		uibridge.OnceUpdateCancelled(func() {
			cancel()
		})
		err := update.Update(ctx, version)
		if err != nil {
			if err.Error() == "cancelled" {
				return false
			}
			uibridge.Endpoint.OnUpdateError(err.Error())
			return true
		}

		uibridge.Endpoint.OnUpdated()
		return true
	}

	return false
}
