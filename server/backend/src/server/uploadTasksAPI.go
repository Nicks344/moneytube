package server

import (
	"net/http"
	"strconv"

	"github.com/Nicks344/moneytube/moneytubemodel"
	"github.com/Nicks344/moneytube/server/backend/src/model"

	"github.com/gin-gonic/gin"
)

func serveUploadTasksAPI(group *gin.RouterGroup) {
	tasksGroup := group.Group("/uploadTasks")
	{
		tasksGroup.GET("/", func(c *gin.Context) {
			key := c.GetHeader("Key")

			tasks, err := model.GetUploadTasks(key)
			if err != nil {
				fail(c, http.StatusInternalServerError, err)
				return
			}

			success(c, tasks)
		})

		tasksGroup.POST("/", func(c *gin.Context) {
			key := c.GetHeader("Key")

			var data moneytubemodel.UploadTask
			var err error
			var id int

			if err = c.BindJSON(&data); err != nil {
				fail(c, http.StatusBadRequest, err)
				return
			}

			if id, err = model.SaveUploadTask(key, data); err != nil {
				fail(c, http.StatusInternalServerError, err)
				return
			}

			success(c, id)
		})

		tasksGroup.POST("/stop/all", func(c *gin.Context) {
			key := c.GetHeader("Key")

			err := model.StopAllTasks(key)
			if err != nil {
				fail(c, http.StatusInternalServerError, err)
				return
			}

			success(c, true)
		})

		tasksGroup.DELETE("/:id", func(c *gin.Context) {
			key := c.GetHeader("Key")
			idStr := c.Param("id")
			if idStr == "all" {
				if err := model.DeleteAllUploadTask(key); err != nil {
					fail(c, http.StatusInternalServerError, err)
					return
				}
			} else {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					fail(c, http.StatusBadRequest, err)
				}

				if err := model.DeleteUploadTask(key, id); err != nil {
					fail(c, http.StatusInternalServerError, err)
					return
				}
			}

			success(c, nil)
		})
	}
}
