package goober

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func (ctx *Context) SendJSON(statusCode int, data interface{}) {
	ctx.Response.status = statusCode
	ctx.Response.Header().Set("Content-Type", "application/json")
	ctx.Response.WriteHeader(statusCode)
	json.NewEncoder(ctx.Response).Encode(data)
}

func (ctx *Context) SendText(statusCode int, text string) {
	ctx.Response.status = statusCode
	ctx.Response.Header().Set("Content-Type", "text/plain")
	ctx.Response.WriteHeader(statusCode)
	ctx.Response.Write([]byte(text))
}

func (ctx *Context) SendFile(filePath string, contentType string) error {
	fmt.Println(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(ctx.Response, "File not found", http.StatusNotFound)
		return err
	}
	defer file.Close()

	ctx.Response.Header().Set("Content-Type", contentType)
	ctx.Response.WriteHeader(http.StatusOK)
	_, err = io.Copy(ctx.Response, file)
	return err
}

func (ctx *Context) SendStream(contentType string, streamData io.Reader) error {
	ctx.Response.Header().Set("Content-Type", contentType)
	ctx.Response.WriteHeader(http.StatusOK)
	_, err := io.Copy(ctx.Response, streamData)
	return err
}

func (ctx *Context) SendError(statusCode int, message string) {
	http.Error(ctx.Response, message, statusCode)
}

func (ctx *Context) Redirect(statusCode int, url string) {
	http.Redirect(ctx.Response, ctx.Request, url, statusCode)
}
