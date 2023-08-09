package update

import (
	"context"
	"errors"
	"os"
	"os/exec"

	"github.com/Nicks344/moneytube/client/core/src/config"
	"github.com/Nicks344/moneytube/client/core/src/serverAPI"

	"github.com/meandrewdev/logger"
)

const appFilename = "resources/app.asar"
const tempFilename = "resources/app.part"

var progress int

func GetProgress() int {
	return progress
}

func Check() (version string, err error) {
	version, err = getVersion()
	if err != nil {
		err = errors.New("Ошибка соединения с сервером обновления")
		return
	}
	if version == config.GetVersion() || version == "0.0" {
		version = ""
	}
	return
}

func Update(ctx context.Context, version string) (err error) {
	logger.NoticeF("new version(%s) downloading", version)
	err = downloadApp(ctx)
	if err != nil {
		if err.Error() != "cancelled" {
			err = errors.New("Ошибка скачивания новой версии")
		}
		return
	}
	logger.Notice("new version downloaded, update")
	err = exec.Command("./update.exe", version).Start()
	if err != nil {
		logger.Error(err)
		err = errors.New("Ошибка старта обновления")
		return
	}
	return
}

func getVersion() (string, error) {
	return serverAPI.GetVersion()
}

func downloadApp(ctx context.Context) error {
	os.Remove(tempFilename)
	err := serverAPI.DownloadUpdate(ctx, tempFilename, &progress)
	progress = 0
	return err

}
