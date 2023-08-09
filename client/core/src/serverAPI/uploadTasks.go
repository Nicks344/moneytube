package serverAPI

import (
	"fmt"

	"github.com/Nicks344/moneytube/moneytubemodel"

	"github.com/imroc/req"
)

func GetUploadTasks() (result []moneytubemodel.UploadTask, err error) {
	var resp *req.Resp
	resp, err = req.Get(host+"/api/user/v1/uploadTasks/", getAuthHeaders())
	if err != nil {
		return
	}

	if err = checkError(resp); err != nil {
		return
	}

	err = parseAnswer(resp, &result)
	return
}

func StopAllTasks() (err error) {
	var resp *req.Resp
	resp, err = req.Post(host+"/api/user/v1/uploadTasks/stop/all", getAuthHeaders())
	if err != nil {
		return
	}

	return checkError(resp)
}

func SaveUploadTask(task moneytubemodel.UploadTask) (result int, err error) {
	var resp *req.Resp
	resp, err = req.Post(host+"/api/user/v1/uploadTasks/", getAuthHeaders(), req.BodyJSON(task))
	if err != nil {
		return
	}

	if err = checkError(resp); err != nil {
		return
	}

	err = parseAnswer(resp, &result)
	return
}

func DeleteUploadTask(id int) (err error) {
	var resp *req.Resp
	resp, err = req.Delete(fmt.Sprintf("%s/api/user/v1/uploadTasks/%d/", host, id), getAuthHeaders())
	if err != nil {
		return
	}

	return checkError(resp)
}

func DeleteAllUploadTasks() (err error) {
	var resp *req.Resp
	resp, err = req.Delete(fmt.Sprintf("%s/api/user/v1/uploadTasks/all/", host), getAuthHeaders())
	if err != nil {
		return
	}

	return checkError(resp)
}
