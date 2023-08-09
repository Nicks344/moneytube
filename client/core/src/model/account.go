package model

import (
	"sync"

	"github.com/Nicks344/moneytube/client/core/src/serverAPI"
	"github.com/Nicks344/moneytube/moneytubemodel"
)

var accountsCache map[int]moneytubemodel.Account
var accountsCacheLock = sync.Mutex{}

func GetAccounts() (result []moneytubemodel.Account, err error) {
	accountsCacheLock.Lock()
	defer accountsCacheLock.Unlock()

	result, err = serverAPI.GetAccounts()
	if err != nil {
		return
	}

	accountsCache = map[int]moneytubemodel.Account{}
	for _, acc := range result {
		accountsCache[acc.ID] = acc
	}

	return
}

func GetAccountsByIDs(ids []int) (result []moneytubemodel.Account, err error) {
	if accountsCache == nil {
		_, err = GetAccounts()
		if err != nil {
			return
		}
	}

	accountsCacheLock.Lock()
	defer accountsCacheLock.Unlock()

	idsMap := map[int]bool{}
	for _, id := range ids {
		idsMap[id] = false
	}

	for id, acc := range accountsCache {
		if _, ok := idsMap[id]; ok {
			result = append(result, acc)
		}
	}

	return
}

func GetAccountsByStatus(status int) (result []moneytubemodel.Account, err error) {
	if accountsCache == nil {
		_, err = GetAccounts()
		if err != nil {
			return
		}
	}

	accountsCacheLock.Lock()
	defer accountsCacheLock.Unlock()

	for _, acc := range accountsCache {
		if acc.Status == status {
			result = append(result, acc)
		}
	}

	return
}

func GetAccount(id int) (result moneytubemodel.Account, err error) {
	if accountsCache == nil {
		_, err = GetAccounts()
		if err != nil {
			return
		}
	}

	accountsCacheLock.Lock()
	defer accountsCacheLock.Unlock()

	result = accountsCache[id]

	return
}

func SaveAccount(account *moneytubemodel.Account) (err error) {
	if accountsCache == nil {
		_, err = GetAccounts()
		if err != nil {
			return
		}
	}

	accountsCacheLock.Lock()
	defer accountsCacheLock.Unlock()

	var id int
	id, err = serverAPI.SaveAccount(*account)
	if err != nil {
		return
	}

	account.ID = id
	accountsCache[id] = *account

	return
}

func DeleteAccount(id int) (err error) {
	accountsCacheLock.Lock()
	defer accountsCacheLock.Unlock()

	delete(accountsCache, id)
	return serverAPI.DeleteAccount(id)
}

func DeleteGroup(group string) (err error) {
	accountsCacheLock.Lock()
	defer accountsCacheLock.Unlock()

	if err := serverAPI.DeleteGroup(group); err != nil {
		return err
	}

	for id, acc := range accountsCache {
		if acc.Group == group {
			acc.Group = ""
			accountsCache[id] = acc
		}
	}

	return
}
