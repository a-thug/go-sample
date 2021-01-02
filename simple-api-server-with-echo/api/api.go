package api
import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func Start(secret string) {
	e := echo.New()

	e.Use(middleware.logger())
	e.Use(middleware.Recover())

	// Allows requests from any origin wth GET, HEAD, PUT, POST or DELETE method
	e.Use(middleware.CORS())

	e.validator = &customValidator{validator: validator.New()}

	e.Logger.SetLevel(log.INFO)

	users := newUserController(secret)

	// Public
	e.GET("/", root)
	e.POST("/api/users", users.Create)

	// Restricted
	r := e.Group("/api", middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(secret)
	}))

	r.GET("/users/:id", users.Get)
	r.PATCH("users/:id", users.Update)
	r.DELETE("/users/:id", users.Delete)

	// Wait for graceful shutdown with 10 sec
	go func() {
		if err := e.Start(":1323"); err != nil {
			e.Logger.Info("Shutting down")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	// Block until an OS signal is received
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	e.Logger.Info("shutdown")
}

func root(c echo.Context) error {
	return c.String(http.StatusOK, "This is root.")
}

type customeValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
