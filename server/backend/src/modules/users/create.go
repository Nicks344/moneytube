package users

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/Nicks344/moneytube/server/backend/src/model"
)

func CreateUser(name string, days int, daysReactivate int, version string) (user model.User, err error) {
	user = model.User{
		Key:            generateKey(name),
		Name:           name,
		Days:           days,
		CreatedAt:      bod(time.Now()),
		DaysReactivate: daysReactivate,
		Version:        version,
	}

	err = model.SaveUser(user)
	if err != nil {
		return
	}

	model.ConnectUser(user.Key)
	return user, err
}

func generateKey(name string) string {
	sha := sha1.New()
	str := fmt.Sprintf("%smegasaltstringa%d", name, time.Now().UnixNano())
	sha.Write([]byte(str))
	return fmt.Sprintf("%x", sha.Sum(nil))
}
