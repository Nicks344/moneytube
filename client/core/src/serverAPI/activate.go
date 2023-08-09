package serverAPI

import (
	"github.com/meandrewdev/logger"
	"github.com/Nicks344/moneytube/client/core/src/config"
	"github.com/Nicks344/moneytube/client/core/src/license/enigma"

	"github.com/imroc/req"
)

type ActivateData struct {
	Key  string
	Name string
}

func Activate(key, hwid string) (data ActivateData, err error) {
	resp, err := req.Post(host+"/api/public/activate", req.BodyJSON(map[string]string{
		"key":  key,
		"hwid": hwid,
	}))
	if err != nil {
		return
	}

	err = checkError(resp)
	if err != nil {
		return
	}

	err = parseAnswer(resp, &data)
	return
}

func Check() bool {
	resp, err := req.Get(host+"/api/user/v1/check", getAuthHeaders())
	if err != nil {
		logger.Warning("License: Check failed due a connect error")
		return false
	}

	if resp.Response().StatusCode != 200 {
		hwid, err := enigma.GetHWID()
		if err != nil {
			logger.Error(err)
			return false
		}
		logger.WarningF("License: Check failed due a check error. Info: API-key: %s, Enigma key: %s, HWID: %s", config.GetApiKey(), config.GetLicenseKey(), hwid)
		return false
	}

	return true
}
