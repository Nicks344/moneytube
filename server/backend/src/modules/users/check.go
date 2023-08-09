package users

import (
	"errors"

	"github.com/Nicks344/moneytube/server/backend/src/model"
)

func Check(key string) (err error) {
	var user model.User
	user, err = model.GetUser(key)
	if err != nil {
		err = errors.New("not exists")
		return
	}

	if !user.IsActivated {
		err = errors.New("not activated")
		return
	}

	if user.DaysLeft < 0 {
		err = errors.New("expired")
		return
	}

	return
}
