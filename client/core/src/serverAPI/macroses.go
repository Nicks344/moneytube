package serverAPI

import (
	"github.com/Nicks344/moneytube/moneytubemodel"

	"github.com/imroc/req"
)

func ExecuteUserMacroses(text string) (result string, err error) {
	result = text
	var resp *req.Resp
	resp, err = req.Post(host+"/api/user/v1/macroses/execute", text, getAuthHeaders())
	if err != nil {
		return
	}

	if err = checkError(resp); err != nil {
		return
	}

	if err = parseAnswer(resp, &result); err != nil {
		result = text
	}

	return
}

func GetMacroses() (result []moneytubemodel.Macros, err error) {
	var resp *req.Resp
	resp, err = req.Get(host+"/api/user/v1/macroses/", getAuthHeaders())
	if err != nil {
		return
	}

	if err = checkError(resp); err != nil {
		return
	}

	err = parseAnswer(resp, &result)
	return
}

func GetMacros(name string) (result moneytubemodel.Macros, err error) {
	var resp *req.Resp
	resp, err = req.Get(host+"/api/user/v1/macroses/"+name, getAuthHeaders())
	if err != nil {
		return
	}

	if err = checkError(resp); err != nil {
		return
	}

	err = parseAnswer(resp, &result)
	return
}

func SaveMacros(macros moneytubemodel.Macros) (err error) {
	var resp *req.Resp
	resp, err = req.Post(host+"/api/user/v1/macroses/", getAuthHeaders(), req.BodyJSON(macros))
	if err != nil {
		return
	}

	if err = checkError(resp); err != nil {
		return
	}

	return
}

func DeleteMacros(name string) (err error) {
	var resp *req.Resp
	resp, err = req.Delete(host+"/api/user/v1/macroses/"+name, getAuthHeaders())
	if err != nil {
		return
	}

	err = checkError(resp)
	return
}
