package main

import (
	"os"

	"github.com/a-thug/go-sample/simple-api-server-with-echo/api"
)

func main() {
	secret := os.Getenv("SECRET")
	if secret == "" {
		panic("missing SECRET")
	}

	api.Start(secret)
}
