package goober

import (
	"fmt"
	"net/http"
)

type Route struct {
	Method     string
	Pattern    string
	Handler    HandlerFunc
	Middleware []MiddlewareFunc
}

type Goober struct {
	*Group
	Logger          bool
	CustomLogStatus bool
	CustomLog       MiddlewareFunc
	groups          []*Group
	middleware      []MiddlewareFunc
}

func Init() *Goober {
	web := &Goober{
		Group: &Group{
			Prefix: "",
		},
		Logger: true,
	}
	return web
}

func Default() *Goober {
	web := &Goober{
		Group: &Group{
			Prefix: "",
		},
		Logger: true,
	}
	return web
}

func (ws *Goober) AddGroup(prefix string) *Group {
	group := &Group{Prefix: prefix}
	ws.groups = append(ws.groups, group)
	return group
}
func (ws *Goober) SetCustomLogger(logger MiddlewareFunc) MiddlewareFunc {
	ws.CustomLog = logger
	return logger
}
func (ws *Goober) Start(port string) error {
	if ws.Logger {
		if ws.CustomLog != nil {
			ws.Use(ws.CustomLog)
		}
		ws.Use(LoggingMiddleware)
	}
	for _, route := range ws.Group.Routes {
		register(ws.Group, &route, ws)
	}
	for _, group := range ws.groups {
		fmt.Println(group)
		for _, route := range group.Routes {
			register(group, &route, ws)
		}
	}
	fmt.Println("Server started on port " + port)
	return http.ListenAndServe(":"+port, nil)
}

func register(group *Group, route *Route, ws *Goober) {
	fullPattern := group.Prefix + route.Pattern
	// fmt.Printf("Registering route: %s %s\n", route.Method, fullPattern)
	http.HandleFunc(fullPattern, func(w http.ResponseWriter, r *http.Request) {
		// middle := route.Middleware
		if r.Method != route.Method {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		ctx := &Context{
			Request: r,
			Response: ResponseWriter{
				ResponseWriter: w,
			},
			Params:  extractParams(fullPattern, r.URL.Path),
			Headers: extractHeaders(r),
			Query:   extractQueryParams(r),
			Body:    extractBody(r),
		}
		//
		allMiddleware := append(ws.middleware, group.Middleware...)
		allMiddleware = append(allMiddleware, route.Middleware...)
		allMiddleware = append(allMiddleware, ws.Group.Middleware...)
		currentHandler := route.Handler
		for i := len(allMiddleware) - 1; i >= 0; i-- {
			next := currentHandler
			mw := allMiddleware[i]
			currentHandler = func(c *Context) {
				mw(c, func() {
					next(c)
				})
			}
		}
		currentHandler(ctx)
	})
}
