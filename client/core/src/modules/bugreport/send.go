package bugreport

import (
	"encoding/json"

	"github.com/Nicks344/moneytube/client/core/src/serverAPI"
)

const (
	ETAEError = "ae-error"
)

type ErrorData struct {
	Error string `json:"error"`
}

func Send(errType string, description string, dataJSON string) (string, error) {
	var errData ErrorData
	err := json.Unmarshal([]byte(dataJSON), &errData)
	if err != nil {
		return "", err
	}

	var data []byte
	switch errType {
	case ETAEError:
		data, err = getAEReport(description, dataJSON)
		if err != nil {
			return "", err
		}
	}

	id, err := serverAPI.SendBugReport(errData.Error, description, data)
	if err != nil {
		return "", err
	}

	return id, nil
}
