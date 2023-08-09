package server

import (
	"net/http"

	"github.com/Nicks344/moneytube/moneytubemodel"
	"github.com/Nicks344/moneytube/server/backend/src/model"

	"github.com/gin-gonic/gin"
)

func serveUploadDataTemplatesAPI(group *gin.RouterGroup) {
	uploadDataTemplatesGroup := group.Group("/uploadDataTemplates")
	{
		uploadDataTemplatesGroup.GET("/", func(c *gin.Context) {
			key := c.GetHeader("Key")

			uploadDataTemplates, err := model.GetUploadDataTemplates(key)
			if err != nil {
				fail(c, http.StatusInternalServerError, err)
				return
			}

			success(c, uploadDataTemplates)
		})

		uploadDataTemplatesGroup.POST("/", func(c *gin.Context) {
			key := c.GetHeader("Key")

			var data moneytubemodel.UploadDataTemplate

			if err := c.BindJSON(&data); err != nil {
				fail(c, http.StatusBadRequest, err)
				return
			}

			err := model.SaveUploadDataTemplate(key, data)
			if err != nil {
				fail(c, http.StatusInternalServerError, err)
				return
			}

			success(c, true)
		})

		uploadDataTemplatesGroup.DELETE("/:id", func(c *gin.Context) {
			key := c.GetHeader("Key")
			id := c.Param("id")

			if err := model.DeleteUploadDataTemplate(key, id); err != nil {
				fail(c, http.StatusInternalServerError, err)
				return
			}

			success(c, nil)
		})
	}
}
