package server

import (
	"net/http"
	"strconv"

	"github.com/Nicks344/moneytube/moneytubemodel"
	"github.com/Nicks344/moneytube/server/backend/src/model"

	"github.com/gin-gonic/gin"
)

func serveUploadDatasAPI(group *gin.RouterGroup) {
	uploadDatasGroup := group.Group("/uploadDatas")
	{
		uploadDatasGroup.GET("/", func(c *gin.Context) {
			key := c.GetHeader("Key")

			uploadDatas, err := model.GetUploadDatas(key)
			if err != nil {
				fail(c, http.StatusInternalServerError, err)
				return
			}

			success(c, uploadDatas)
		})

		uploadDatasGroup.POST("/", func(c *gin.Context) {
			key := c.GetHeader("Key")

			var data moneytubemodel.UploadData
			var err error
			var answer struct {
				ID    int
				Tasks []moneytubemodel.UploadTask
			}

			if err = c.BindJSON(&data); err != nil {
				fail(c, http.StatusBadRequest, err)
				return
			}

			if answer.ID, answer.Tasks, err = model.SaveUploadData(key, data); err != nil {
				fail(c, http.StatusInternalServerError, err)
				return
			}

			success(c, answer)
		})

		uploadDatasGroup.DELETE("/:id", func(c *gin.Context) {
			key := c.GetHeader("Key")
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				fail(c, http.StatusBadRequest, err)
			}

			if err := model.DeleteUploadData(key, id); err != nil {
				fail(c, http.StatusInternalServerError, err)
				return
			}

			success(c, nil)
		})
	}
}
