package main

import (
	"fmt"
	"runtime/debug"

	"github.com/a-thug/go-sample/simple-api-server-with-gin/api"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
		}
	}()
	r := api.SetupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
