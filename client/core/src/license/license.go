package license

import (
	"errors"

	"github.com/meandrewdev/logger"
	"github.com/Nicks344/moneytube/client/core/src/config"
	"github.com/Nicks344/moneytube/client/core/src/license/enigma"
	"github.com/Nicks344/moneytube/client/core/src/serverAPI"
)

func Check() bool {
	key := config.GetApiKey()
	if key == "" {
		logger.Warning("License: API key is empty")
		return false
	}

	if !enigma.LoadAndCheckKey() {
		logger.Warning("License: Cannot load and check enigma key")
		return false
	}

	if enigma.GetDaysLeft() < 0 {
		logger.Warning("License: No more days of use")
		return false
	}

	return serverAPI.Check()
}

func Register(key string) error {
	if err := config.SetApiKey(key); err != nil {
		logger.Error(err)
		return errors.New("cannot save api key")
	}

	hwid, err := enigma.GetHWID()
	if err != nil {
		logger.Error(err)
		return errors.New("cannot get hardware id")
	}
	answer, err := serverAPI.Activate(key, hwid)
	if err != nil {
		logger.Error(err)
		return errors.New("activate error")
	}

	if !enigma.CheckAndSaveKey(answer.Name, answer.Key) {
		return errors.New("invalid enigma key")
	}

	if err := config.SetLicenseKey(answer.Key); err != nil {
		logger.Error(err)
		return errors.New("cannot save license key")
	}

	return nil
}
