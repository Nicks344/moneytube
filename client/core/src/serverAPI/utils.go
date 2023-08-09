package serverAPI

import (
	"errors"

	"github.com/meandrewdev/logger"
	"github.com/Nicks344/moneytube/client/core/src/config"
	"github.com/Nicks344/moneytube/client/core/src/license/enigma"
	"github.com/Nicks344/moneytube/licensehash"

	"github.com/mitchellh/mapstructure"

	"github.com/imroc/req"
)

var host string

type answer struct {
	Error  string
	Result interface{}
}

func getAuthHeaders() (headers req.Header) {
	headers = req.Header{
		"Version": "2",
	}

	hwid, err := enigma.GetHWID()
	if err != nil {
		logger.Error(err)
		return
	}

	key := config.GetLicenseKey()

	headers["Hash"] = licensehash.GetInfoHash(key, hwid)
	headers["Key"] = config.GetApiKey()

	return
}

func checkError(resp *req.Resp) error {
	if resp.Response().StatusCode != 200 {
		return parseAnswer(resp, nil)
	}

	return nil
}

func parseAnswer(resp *req.Resp, result interface{}) error {
	var ans answer
	err := resp.ToJSON(&ans)
	if err != nil {
		return err
	}

	if ans.Error != "" {
		return errors.New(ans.Error)
	}

	return mapstructure.Decode(ans.Result, &result)
}
