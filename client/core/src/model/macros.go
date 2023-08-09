package model

import (
	"sync"

	"github.com/Nicks344/moneytube/client/core/src/serverAPI"
	"github.com/Nicks344/moneytube/moneytubemodel"
)

var macrosesCache map[string]moneytubemodel.Macros
var macrosesCacheLock = sync.Mutex{}

func GetMacroses() (result []moneytubemodel.Macros, err error) {
	macrosesCacheLock.Lock()
	defer macrosesCacheLock.Unlock()

	result, err = serverAPI.GetMacroses()
	if err != nil {
		return
	}

	macrosesCache = map[string]moneytubemodel.Macros{}
	for _, macros := range result {
		macrosesCache[macros.Name] = macros
	}

	return
}

func GetMacros(name string) (result moneytubemodel.Macros, err error) {
	if macrosesCache == nil {
		_, err = GetMacroses()
		if err != nil {
			return
		}
	}

	macrosesCacheLock.Lock()
	defer macrosesCacheLock.Unlock()

	result = macrosesCache[name]

	return
}

func SaveMacros(macros moneytubemodel.Macros) (err error) {
	if macrosesCache == nil {
		_, err = GetMacroses()
		if err != nil {
			return
		}
	}

	macrosesCacheLock.Lock()
	defer macrosesCacheLock.Unlock()

	macrosesCache[macros.Name] = macros

	return serverAPI.SaveMacros(macros)
}

func DeleteMacros(name string) (err error) {
	macrosesCacheLock.Lock()
	defer macrosesCacheLock.Unlock()

	delete(macrosesCache, name)
	return serverAPI.DeleteMacros(name)
}
