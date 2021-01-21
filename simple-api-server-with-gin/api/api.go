package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// root
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "This is root.")
	})

	return r
}
