package goober

import (
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

type Group struct {
	Prefix     string
	Routes     []Route
	Middleware []MiddlewareFunc
}

func NewGroup(prefix string) *Group {
	return &Group{Prefix: prefix}
}

func (g *Group) GET(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.addRoute("GET", pattern, handler, middleware...)
}

func (g *Group) Use(middleware ...MiddlewareFunc) bool {
	g.Middleware = append(g.Middleware, middleware[0])
	return true
}

func (g *Group) POST(pattern string, handler HandlerFunc) {
	g.addRoute("POST", pattern, handler)
}

func (g *Group) DELETE(pattern string, handler HandlerFunc) {
	g.addRoute("DELETE", pattern, handler)
}

func (g *Group) PUT(pattern string, handler HandlerFunc) {
	g.addRoute("PUT", pattern, handler)
}

func (g *Group) PATCH(pattern string, handler HandlerFunc) {
	g.addRoute("PATCH", pattern, handler)
}

func (g *Group) addRoute(method, pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.Routes = append(g.Routes, Route{Method: method, Pattern: pattern, Handler: handler, Middleware: middleware})
}

func (g *Group) ServeStatic(urlPrefix, directory string) {
	if !strings.HasSuffix(urlPrefix, "/") {
		urlPrefix += "/"
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range entries {
		fmt.Println(e.Name())
		g.Routes = append(g.Routes, Route{
			Method:  "GET",
			Pattern: urlPrefix + e.Name(),
			Handler: func(ctx *Context) {
				ctx.SendFile(directory+"/"+e.Name(), g.getContentType(e.Name()))
			},
		})
	}
}

func (g *Group) getContentType(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return "application/octet-stream"
	}
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		return "application/octet-stream"
	}
	return contentType
}
