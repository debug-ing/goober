package goober

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func extractParams(pattern, path string) map[string]string {
	params := make(map[string]string)
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	for i, part := range patternParts {
		if strings.HasPrefix(part, ":") && i < len(pathParts) {
			key := strings.TrimPrefix(part, ":")
			params[key] = pathParts[i]
		}
	}

	return params
}

func extractHeaders(r *http.Request) map[string]string {
	headers := make(map[string]string)
	for name, values := range r.Header {
		headers[name] = values[0]
	}
	return headers
}

func extractQueryParams(r *http.Request) map[string]string {
	query := make(map[string]string)
	for name, values := range r.URL.Query() {
		query[name] = values[0]
	}
	return query
}

func extractBody(r *http.Request) map[string]interface{} {
	body := make(map[string]interface{})
	if r.Body == nil {
		return body
	}
	defer r.Body.Close()

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil || len(bodyBytes) == 0 {
		return body
	}
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		return body
	}
	return body
}

func extractCookies(r *http.Request) map[string]string {
	cookies := make(map[string]string)
	for _, cookie := range r.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}
	return cookies
}
