package model

import (
	"errors"
	"sync"

	"github.com/Nicks344/moneytube/client/core/src/serverAPI"
	"github.com/Nicks344/moneytube/moneytubemodel"
)

var uploadDataTemplatesCache map[string]moneytubemodel.UploadDataTemplate
var uploadDataTemplatesCacheLock = sync.Mutex{}

func GetUploadDataTemplates() (result []moneytubemodel.UploadDataTemplate, err error) {
	uploadDataTemplatesCacheLock.Lock()
	defer uploadDataTemplatesCacheLock.Unlock()

	result, err = serverAPI.GetUploadDataTemplates()
	if err != nil {
		return
	}

	uploadDataTemplatesCache = map[string]moneytubemodel.UploadDataTemplate{}
	for _, templ := range result {
		uploadDataTemplatesCache[templ.Label] = templ
	}

	return
}

func GetUploadDataTemplate(id string) (result moneytubemodel.UploadDataTemplate, err error) {
	if uploadDataTemplatesCache == nil {
		_, err = GetUploadDataTemplates()
		if err != nil {
			return
		}
	}

	uploadDataTemplatesCacheLock.Lock()
	defer uploadDataTemplatesCacheLock.Unlock()

	result = uploadDataTemplatesCache[id]

	return
}

func GetUploadDataTemplateByName(name string) (result moneytubemodel.UploadDataTemplate, err error) {
	if uploadDataTemplatesCache == nil {
		_, err = GetUploadDataTemplates()
		if err != nil {
			return
		}
	}

	uploadDataTemplatesCacheLock.Lock()
	defer uploadDataTemplatesCacheLock.Unlock()

	for _, result = range uploadDataTemplatesCache {
		if result.Label == name {
			return
		}
	}

	result = moneytubemodel.UploadDataTemplate{}
	err = errors.New("upload template not found")
	return
}

func SaveUploadDataTemplate(uploadDataTemplate *moneytubemodel.UploadDataTemplate) (err error) {
	if uploadDataTemplatesCache == nil {
		_, err = GetUploadDataTemplates()
		if err != nil {
			return
		}
	}

	uploadDataTemplatesCacheLock.Lock()
	defer uploadDataTemplatesCacheLock.Unlock()

	err = serverAPI.SaveUploadDataTemplate(*uploadDataTemplate)
	if err != nil {
		return
	}

	uploadDataTemplatesCache[uploadDataTemplate.Label] = *uploadDataTemplate

	return
}

func DeleteUploadDataTemplate(id string) (err error) {
	uploadDataTemplatesCacheLock.Lock()
	defer uploadDataTemplatesCacheLock.Unlock()

	delete(uploadDataTemplatesCache, id)
	return serverAPI.DeleteUploadDataTemplate(id)
}
