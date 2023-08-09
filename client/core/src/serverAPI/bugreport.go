package serverAPI

import (
	"encoding/base64"

	"github.com/imroc/req"
)

type ReportInfo struct {
	Error       string `json:"error"`
	Description string `json:"description"`
	Data        string `json:"data"`
}

func SendBugReport(errText string, desc string, archive []byte) (result string, err error) {
	data := base64.StdEncoding.EncodeToString(archive)
	report := ReportInfo{
		Error:       errText,
		Description: desc,
		Data:        data,
	}
	var resp *req.Resp
	resp, err = req.Post(host+"/api/user/v1/bugreport/", getAuthHeaders(), req.BodyJSON(report))
	if err != nil {
		return
	}

	if err = checkError(resp); err != nil {
		return
	}

	err = parseAnswer(resp, &result)
	return
}
