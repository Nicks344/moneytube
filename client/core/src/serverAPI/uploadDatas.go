package serverAPI

import (
	"fmt"

	"github.com/Nicks344/moneytube/moneytubemodel"

	"github.com/imroc/req"
)

func GetUploadDatas() (result []moneytubemodel.UploadData, err error) {
	var resp *req.Resp
	resp, err = req.Get(host+"/api/user/v1/uploadDatas/", getAuthHeaders())
	if err != nil {
		return
	}

	if err = checkError(resp); err != nil {
		return
	}

	err = parseAnswer(resp, &result)
	return
}

func SaveUploadData(task moneytubemodel.UploadData) (id int, tasks []moneytubemodel.UploadTask, err error) {
	var result struct {
		ID    int
		Tasks []moneytubemodel.UploadTask
	}

	var resp *req.Resp
	resp, err = req.Post(host+"/api/user/v1/uploadDatas/", getAuthHeaders(), req.BodyJSON(task))
	if err != nil {
		return
	}

	if err = checkError(resp); err != nil {
		return
	}

	if err = parseAnswer(resp, &result); err != nil {
		return
	}

	id = result.ID
	tasks = result.Tasks

	return
}

func DeleteUploadData(id int) (err error) {
	var resp *req.Resp
	resp, err = req.Delete(fmt.Sprintf("%s/api/user/v1/uploadDatas/%d/", host, id), getAuthHeaders())
	if err != nil {
		return
	}

	if err = checkError(resp); err != nil {
		return
	}

	return
}
