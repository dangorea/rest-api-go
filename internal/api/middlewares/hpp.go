package middlewares

import (
	"net/http"
	"strings"
)

type HPPOptions struct {
	CheckQuery                  bool
	CheckBody                   bool
	CheckBodyOnlyForContentType string
	Whitelist                   []string
}

func Hpp(options HPPOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if options.CheckBody && r.Method == http.MethodPost && isCorrectContentType(r, options.CheckBodyOnlyForContentType) {
				// filter the body parameters
				filterBodyParams(r, options.Whitelist)
			}

			if options.CheckQuery && r.URL.Query() != nil {
				// filter the query parameters
				filterQueryParams(r, options.Whitelist)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isCorrectContentType(r *http.Request, contentType string) bool {
	return strings.Contains(r.Header.Get("Content-Type"), contentType)
}

func filterBodyParams(r *http.Request, whitelist []string) {
	err := r.ParseForm()
	if err != nil {
		return
	}

	for key, values := range r.Form {
		if len(values) > 1 {
			r.Form.Set(key, values[0]) // First value
			// r.Form.Set(key, values[len(values)-1]) // Accept the last value
		}

		if !isWhitelisted(key, whitelist) {
			delete(r.Form, key)
		}
	}
}

func filterQueryParams(r *http.Request, whitelist []string) {
	query := r.URL.Query()

	for key, values := range query {
		if len(values) > 1 {
			query.Set(key, values[0]) // First value
			// query.Set(key, values[len(values)-1]) // Accept the last value
		}

		if !isWhitelisted(key, whitelist) {
			query.Del(key)
		}
	}

	r.URL.RawQuery = query.Encode()
}

func isWhitelisted(param string, whitelist []string) bool {
	for _, v := range whitelist {
		if param == v {
			return true
		}
	}

	return false
}
