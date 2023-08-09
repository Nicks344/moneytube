package server

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Nicks344/moneytube/server/backend/src/model"

	"github.com/gin-gonic/gin"
)

func serveBugReportAPI(group *gin.RouterGroup) {
	group.POST("/bugreport/", func(c *gin.Context) {
		key := c.GetHeader("Key")

		var input struct {
			Error       string `json:"error"`
			Description string `json:"description"`
			Data        string `json:"data"`
		}

		if err := c.BindJSON(&input); err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}

		user, err := model.GetUser(key)
		if err != nil {
			fail(c, http.StatusInternalServerError, err)
			return
		}

		report := model.Report{
			UserName:    user.Name,
			Error:       input.Error,
			Description: input.Description,
		}

		id, err := model.SaveReport(report)
		if err != nil {
			fail(c, http.StatusInternalServerError, err)
			return
		}

		archive, err := base64.StdEncoding.DecodeString(input.Data)
		if err != nil {
			fail(c, http.StatusInternalServerError, err)
			return
		}

		if err := ioutil.WriteFile(fmt.Sprintf("data/reports/%s.zip", id), archive, 0666); err != nil {
			fail(c, http.StatusInternalServerError, err)
			return
		}

		success(c, id)
	})
}
