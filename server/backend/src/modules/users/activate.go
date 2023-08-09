package users

import (
	"errors"
	"sync"
	"time"

	"github.com/Nicks344/moneytube/server/backend/src/config"
	"github.com/Nicks344/moneytube/server/backend/src/model"
	"github.com/Nicks344/moneytube/server/backend/src/modules/enigmakeygen"
)

var activateLock = sync.Mutex{}
var activateLocks = map[string]*sync.Mutex{}

func Activate(key, hwid string) (enigmaKey string, name string, err error) {
	activateLock.Lock()
	lock, ok := activateLocks[key]
	if !ok {
		lock = &sync.Mutex{}
		activateLocks[key] = lock
	}
	activateLock.Unlock()

	lock.Lock()
	defer func() {
		activateLock.Lock()
		delete(activateLocks, key)
		lock.Unlock()
		activateLock.Unlock()
	}()

	var user model.User
	user, err = model.GetUser(key)
	if err != nil {
		return
	}

	if user.IsActivated {
		if user.HWID == hwid {
			return user.EnigmaKey, user.Name, nil
		}

		daysFromActivate := int(time.Now().Sub(user.ActivatedAt).Hours() / 24)
		if daysFromActivate < user.DaysReactivate {
			err = errors.New("user is already activated")
			return
		}
	}

	name = user.Name

	enigmaKey, err = enigmakeygen.GenerateKey(config.GetEnigmaProject(), name, hwid, user.Days)
	if err != nil {
		return
	}

	user.EnigmaKey = enigmaKey
	user.HWID = hwid
	user.ActivatedAt = bod(time.Now())
	user.IsActivated = true
	user.IsActive = true
	err = model.SaveUser(user)
	model.ConnectUser(user.Key)
	return
}
