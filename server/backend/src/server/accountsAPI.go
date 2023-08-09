package server

import (
	"net/http"
	"strconv"

	"github.com/Nicks344/moneytube/moneytubemodel"
	"github.com/Nicks344/moneytube/server/backend/src/model"

	"github.com/gin-gonic/gin"
)

func serveAccountsAPI(group *gin.RouterGroup) {
	accountsGroup := group.Group("/accounts")
	{
		accountsGroup.GET("/", func(c *gin.Context) {
			key := c.GetHeader("Key")

			accounts, err := model.GetAccounts(key)

			if err != nil {
				fail(c, http.StatusInternalServerError, err)
				return
			}

			success(c, accounts)
		})

		accountsGroup.POST("/", func(c *gin.Context) {
			key := c.GetHeader("Key")

			var account moneytubemodel.Account
			var id int
			var err error

			if err = c.BindJSON(&account); err != nil {
				fail(c, http.StatusBadRequest, err)
				return
			}

			if id, err = model.SaveAccount(key, account); err != nil {
				fail(c, http.StatusInternalServerError, err)
				return
			}

			success(c, id)
		})

		accountsGroup.DELETE("/:id", func(c *gin.Context) {
			key := c.GetHeader("Key")
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				fail(c, http.StatusBadRequest, err)
			}

			if err := model.DeleteAccount(key, id); err != nil {
				fail(c, http.StatusInternalServerError, err)
				return
			}

			success(c, nil)
		})
	}

	groupsGroup := group.Group("/groups")
	{
		groupsGroup.DELETE("/:id", func(c *gin.Context) {
			key := c.GetHeader("Key")
			id := c.Param("id")

			if err := model.DeleteGroup(key, id); err != nil {
				fail(c, http.StatusInternalServerError, err)
				return
			}

			success(c, nil)
		})
	}
}
