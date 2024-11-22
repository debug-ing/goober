package goober

import (
	"fmt"
	"time"
)

type MiddlewareFunc func(ctx *Context, next func())

func LoggingMiddleware(ctx *Context, next func()) {
	start := time.Now()
	url := ctx.Request.URL
	method := ctx.Request.Method
	next()
	end := time.Now()
	let := end.Sub(start)
	currentTime := time.Now().Format("2006/01/02 - 15:04:05")
	fmt.Printf("[GOOBER] %s | %d | %s | %s | %s | %s\n", currentTime, ctx.Response.status, let, ctx.GetIP(), method, url)
}
