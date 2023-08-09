package serverAPI

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/Nicks344/moneytube/moneytubemodel"

	"github.com/imroc/req"
)

func GetUploadDataTemplates() (result []moneytubemodel.UploadDataTemplate, err error) {
	var resp *req.Resp
	resp, err = req.Get(host+"/api/user/v1/uploadDataTemplates/", getAuthHeaders())
	if err != nil {
		return
	}

	if err = checkError(resp); err != nil {
		return
	}

	err = parseAnswer(resp, &result)
	return
}

func SaveUploadDataTemplate(task moneytubemodel.UploadDataTemplate) (err error) {

	str, err := json.Marshal(task.UploadDataFields)
	if err == nil {
		file, err := os.Create(filepath.Join("temp", task.Label+".tmpl"))
		if err == nil {
			defer file.Close()
			file.WriteString(string(str))
		}
	} else {
		file1, _ := os.Create(filepath.Join("temp", "error.tmpl"))
		file1.WriteString(err.Error())
		file1.Close()
	}

	var resp *req.Resp
	resp, err = req.Post(host+"/api/user/v1/uploadDataTemplates/", getAuthHeaders(), req.BodyJSON(task))
	if err != nil {
		return
	}

	if err = checkError(resp); err != nil {
		return
	}

	return
}

func DeleteUploadDataTemplate(id string) (err error) {
	var resp *req.Resp
	id = url.PathEscape(id)
	resp, err = req.Delete(fmt.Sprintf("%s/api/user/v1/uploadDataTemplates/%s/", host, id), getAuthHeaders())
	if err != nil {
		return
	}

	if err = checkError(resp); err != nil {
		return
	}

	return
}
