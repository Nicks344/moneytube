package serverAPI

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/Nicks344/moneytube/client/core/src/utils"

	"github.com/imroc/req"
)

func GetVersion() (result string, err error) {
	var resp *req.Resp
	resp, err = req.Get(host+"/api/update/version", getAuthHeaders())
	if err != nil {
		return
	}

	if err = checkError(resp); err != nil {
		return
	}

	err = parseAnswer(resp, &result)
	return
}

func DownloadUpdate(ctx context.Context, resultFile string, progress *int) (err error) {
	req.SetTimeout(time.Minute * 30)
	resp, err := req.Get(fmt.Sprintf("%s/api/update/download", host), getAuthHeaders())
	if err != nil {
		return err
	}
	defer resp.Response().Body.Close()
	if resp.Response().StatusCode != 200 {
		return errors.New("update download error, status: " + resp.Response().Status)
	}

	os.Remove(resultFile)
	tempFile, err := os.Create(resultFile)
	if err != nil {
		return err
	}
	defer tempFile.Close()
	len, _ := strconv.ParseInt(resp.Response().Header.Get("Content-Length"), 10, 64)
	_, err = io.Copy(tempFile, &passThru{Reader: resp.Response().Body, total: len, progress: progress, ctx: ctx})
	return err
}

type passThru struct {
	io.Reader
	total    int64
	readed   int64
	progress *int
	ctx      context.Context
}

func (pt *passThru) Read(p []byte) (int, error) {
	if utils.IsContextCancelled(pt.ctx) {
		return 0, errors.New("cancelled")
	}
	n, err := pt.Reader.Read(p)
	pt.readed += int64(n)
	*pt.progress = int((float64(pt.readed) / float64(pt.total)) * 100)

	return n, err
}
