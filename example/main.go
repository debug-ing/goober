package main

import (
	"fmt"
	"goober"
	"net/http"
)

func main() {
	server := goober.Init()
	apiGroup := server.AddGroup("/api")
	apiGroup.GET("/hello", func(ctx *goober.Context) {
		id := ctx.Query["id"]
		fmt.Println(id)
		ctx.SendText(200, "HI"+id)
		fmt.Println("HI")
	})
	server.GET("/test", func(ctx *goober.Context) {
		id := ctx.Query["id"]
		fmt.Println(id)
		ctx.SendText(200, "HI"+id)
	})
	server.GET("/redirect", func(ctx *goober.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "https://www.google.com")
	})

	server.AddGroup("/file").ServeStatic("/assets", "./static")
	if err := server.Start("8082"); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
