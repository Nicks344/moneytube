package model

import (
	"sync"

	"github.com/Nicks344/moneytube/client/core/src/serverAPI"
	"github.com/Nicks344/moneytube/moneytubemodel"
)

var uploadTasksCache map[int]moneytubemodel.UploadTask
var uploadTasksCacheLock = sync.Mutex{}

func GetUploadTasks() (result []moneytubemodel.UploadTask, err error) {
	uploadTasksCacheLock.Lock()
	defer uploadTasksCacheLock.Unlock()

	if uploadTasksCache == nil {
		result, err = serverAPI.GetUploadTasks()
		if err != nil {
			return
		}

		uploadTasksCache = map[int]moneytubemodel.UploadTask{}
		for i := range result {
			loadUploadTaskDependencies(&result[i])
			uploadTasksCache[result[i].ID] = result[i]
		}
	} else {
		for _, task := range uploadTasksCache {
			result = append(result, task)
		}
	}

	return
}

func GetUploadTasksByStatuses(statuses []int) (result []moneytubemodel.UploadTask, err error) {
	if uploadTasksCache == nil {
		_, err = GetUploadTasks()
		if err != nil {
			return
		}
	}

	uploadTasksCacheLock.Lock()
	defer uploadTasksCacheLock.Unlock()

	statusesMap := map[int]bool{}
	for _, id := range statuses {
		statusesMap[id] = false
	}

	for _, task := range uploadTasksCache {
		if _, ok := statusesMap[task.Status]; ok {
			result = append(result, task)
		}
	}

	return
}

func GetUploadTask(id int) (result moneytubemodel.UploadTask, err error) {
	if uploadTasksCache == nil {
		_, err = GetUploadTasks()
		if err != nil {
			return
		}
	}

	uploadTasksCacheLock.Lock()
	defer uploadTasksCacheLock.Unlock()

	result = uploadTasksCache[id]

	return
}

func SaveUploadTask(uploadTask *moneytubemodel.UploadTask) (err error) {
	if uploadTasksCache == nil {
		_, err = GetUploadTasks()
		if err != nil {
			return
		}
	}

	uploadTasksCacheLock.Lock()
	defer uploadTasksCacheLock.Unlock()

	var id int
	id, err = serverAPI.SaveUploadTask(*uploadTask)
	if err != nil {
		return
	}

	uploadTask.ID = id
	loadUploadTaskDependencies(uploadTask)
	uploadTasksCache[id] = *uploadTask

	return
}

func DeleteUploadTask(id int) (err error) {
	uploadTasksCacheLock.Lock()
	defer uploadTasksCacheLock.Unlock()

	task := uploadTasksCache[id]
	delete(uploadTasksCache, id)

	var ok bool
	for _, tsk := range uploadTasksCache {
		if tsk.DetailsID == task.DetailsID {
			ok = true
			break
		}
	}

	if !ok {
		err = DeleteUploadData(task.DetailsID)
		if err != nil {
			return
		}
	}

	return serverAPI.DeleteUploadTask(id)
}

func DeleteAllUploadTasks() (err error) {
	uploadTasksCacheLock.Lock()
	defer uploadTasksCacheLock.Unlock()

	uploadDatasCacheLock.Lock()
	defer uploadDatasCacheLock.Unlock()

	uploadTasksCache = map[int]moneytubemodel.UploadTask{}
	uploadDatasCache = map[int]*moneytubemodel.UploadData{}
	return serverAPI.DeleteAllUploadTasks()
}

func DeleteUploadTasksWithAccountID(id int) (err error) {
	uploadTasksCacheLock.Lock()
	defer uploadTasksCacheLock.Unlock()

	for _, task := range uploadTasksCache {
		if task.AccountID == id {
			delete(uploadTasksCache, task.ID)
			err = serverAPI.DeleteUploadTask(task.ID)
			if err != nil {
				return
			}
		}
	}

	return
}

func loadUploadTaskDependencies(task *moneytubemodel.UploadTask) (err error) {
	task.Account, err = GetAccount(task.AccountID)
	if err != nil {
		return
	}

	task.Details, err = GetUploadData(task.DetailsID)
	return
}
