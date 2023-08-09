package server

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Nicks344/moneytube/licensehash"
	"github.com/Nicks344/moneytube/server/backend/src/config"
	"github.com/Nicks344/moneytube/server/backend/src/model"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
)

func Serve(debug bool) {
	router := gin.Default()

	if debug {
		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET,POST,PATCH,PUT,DELETE"},
			AllowHeaders:     []string{"content-type"},
			ExposeHeaders:    []string{"*"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}

	router.Use(version())

	serveAdminAPI(router, debug)

	servePublicAPI(router)

	serveUserAPI(router)

	serveUpdateAPI(router)

	if !debug {
		log.Fatal(autotls.RunWithManager(router, &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(config.GetDomain()),
			Cache:      autocert.DirCache("/var/www/.cache"),
			ForceRSA:   true,
		}))
	}

	log.Fatal(router.Run("127.0.0.1:" + config.GetPort()))
}

func version() gin.HandlerFunc {
	return func(c *gin.Context) {
		v := c.GetHeader("Version")
		c.Set("version", v)
		c.Next()
	}
}

var invalidKeyErr = errors.New("invalid key")

func userAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		version := c.GetString("version")
		key := c.GetHeader("Key")
		hash := c.GetHeader("Hash")

		if key == "" {
			fail(c, 403, invalidKeyErr)
			return
		}

		user, err := model.GetUser(key)
		if err != nil {
			fail(c, 403, invalidKeyErr)
			return
		}

		if !user.IsActivated || !user.IsActive || user.Version != version {
			fail(c, 403, invalidKeyErr)
			return
		}

		if hash != licensehash.GetInfoHash(user.EnigmaKey, user.HWID) {
			fail(c, 403, invalidKeyErr)
			return
		}
	}
}

func success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"result": data,
	})
}

func fail(c *gin.Context, code int, err error) {
	c.JSON(code, map[string]interface{}{
		"error": err.Error(),
	})
	c.Abort()
}
