package serverAPI

import (
	"fmt"
	"net/url"

	"github.com/Nicks344/moneytube/moneytubemodel"

	"github.com/imroc/req"
)

func GetAccounts() (result []moneytubemodel.Account, err error) {
	var resp *req.Resp
	resp, err = req.Get(host+"/api/user/v1/accounts/", getAuthHeaders())
	if err != nil {
		return
	}

	if err = checkError(resp); err != nil {
		return
	}

	err = parseAnswer(resp, &result)
	return
}

func SaveAccount(account moneytubemodel.Account) (result int, err error) {
	var resp *req.Resp
	resp, err = req.Post(host+"/api/user/v1/accounts/", getAuthHeaders(), req.BodyJSON(account))
	if err != nil {
		return
	}

	if err = checkError(resp); err != nil {
		return
	}

	err = parseAnswer(resp, &result)
	return
}

func DeleteAccount(id int) (err error) {
	var resp *req.Resp
	resp, err = req.Delete(fmt.Sprintf("%s/api/user/v1/accounts/%d/", host, id), getAuthHeaders())
	if err != nil {
		return
	}

	if err = checkError(resp); err != nil {
		return
	}

	return
}

func DeleteGroup(id string) (err error) {
	var resp *req.Resp
	id = url.PathEscape(id)
	resp, err = req.Delete(fmt.Sprintf("%s/api/user/v1/groups/%s/", host, id), getAuthHeaders())
	if err != nil {
		return
	}

	if err = checkError(resp); err != nil {
		return
	}

	return
}
