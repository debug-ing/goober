package goober

import "net/http"

type Context struct {
	Request  *http.Request
	Response ResponseWriter
	Params   map[string]string
	Headers  map[string]string
	Query    map[string]string
	Body     map[string]interface{}
	//
	Cookies  map[string]string
	ClientIP string
}

type ResponseWriter struct {
	http.ResponseWriter
	status int
}

type HandlerFunc func(ctx *Context)

func (ctx *Context) GetIP() string {
	ip := ctx.Request.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = ctx.Request.RemoteAddr
	}
	return ip
}
