package api

import (
	"log"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Api  *gin.Engine
	Auth *Auth
}

// SetupRouter is function of setting route
func SetupRouter() *gin.Engine {
	r := Router{
		Api:  gin.Default(),
		Auth: NewAuthService(),
	}

	// Root
	r.Api.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "This is root.")
	})

	// Authorization group
	authorized := r.Api.Group("/api")
	r.makeAuthRoutes(authorized)

	return r.Api
}

func (r *Router) makeAuthRoutes(authorized *gin.RouterGroup) gin.IRoutes {
	mw := r.Auth.getMiddleware()

	r.Api.NoRoute(mw.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		// FIXME: Do not display user's password at log.
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	authorized.Use(mw.MiddlewareFunc())
	authorized.GET("/hello", helloHandler)
	return authorized
}

func helloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)

	c.JSON(200, gin.H{
		"user_id": claims["user_id"],
		"claims":  claims,
		"text":    "Hello World.",
	})
}
