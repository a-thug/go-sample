package main

import (
	"net/http"
	"time"

	"github.com/a-thug/go-sample/simple-api-server-without-framework/api"
	"go.uber.org/zap"
)

func main() {
	logger := zap.NewExample().Sugar()
	defer logger.Sync() // flushes buffer, if any

	router := http.NewServeMux()

	api.Setup(router, logger)

	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		Handler:      router,
	}

	logger.Infow("Start server", "addr", server.Addr)
	logger.Fatal(server.ListenAndServe())
}
