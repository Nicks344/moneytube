package model

import (
	"sync"

	"github.com/Nicks344/moneytube/client/core/src/serverAPI"
	"github.com/Nicks344/moneytube/moneytubemodel"
)

var uploadDatasCache map[int]*moneytubemodel.UploadData
var uploadDatasCacheLock = sync.Mutex{}

func GetUploadDatas() (result []*moneytubemodel.UploadData, err error) {
	uploadDatasCacheLock.Lock()
	defer uploadDatasCacheLock.Unlock()

	datas, err := serverAPI.GetUploadDatas()
	if err != nil {
		return
	}

	result = make([]*moneytubemodel.UploadData, len(datas))
	for i := range datas {
		result[i] = &datas[i]
	}

	uploadDatasCache = map[int]*moneytubemodel.UploadData{}
	for _, acc := range result {
		uploadDatasCache[acc.ID] = acc
	}

	return
}

func GetUploadData(id int) (result *moneytubemodel.UploadData, err error) {
	if uploadDatasCache == nil {
		_, err = GetUploadDatas()
		if err != nil {
			return
		}
	}

	uploadDatasCacheLock.Lock()
	defer uploadDatasCacheLock.Unlock()

	result = uploadDatasCache[id]

	return
}

func SaveUploadData(uploadData *moneytubemodel.UploadData) (tasks []moneytubemodel.UploadTask, err error) {
	if uploadDatasCache == nil {
		_, err = GetUploadDatas()
		if err != nil {
			return
		}
	}

	uploadDatasCacheLock.Lock()

	var id int
	id, tasks, err = serverAPI.SaveUploadData(*uploadData)
	if err != nil {
		uploadDatasCacheLock.Unlock()
		return
	}

	uploadData.ID = id
	uploadDatasCache[id] = uploadData

	if uploadTasksCache == nil {
		_, err = GetUploadTasks()
		if err != nil {
			uploadDatasCacheLock.Unlock()
			return
		}
	}

	uploadDatasCacheLock.Unlock()

	for i := range tasks {
		loadUploadTaskDependencies(&tasks[i])
		uploadTasksCacheLock.Lock()
		uploadTasksCache[tasks[i].ID] = tasks[i]
		uploadTasksCacheLock.Unlock()
	}

	return
}

func DeleteUploadData(id int) (err error) {
	uploadDatasCacheLock.Lock()
	defer uploadDatasCacheLock.Unlock()

	delete(uploadDatasCache, id)
	return serverAPI.DeleteUploadData(id)
}
