package middlewares

import (
	"fmt"
	"net/http"
	"strings"
)

// HPP stands for HTTP Parameter Pollution
// These cleaner functions are used for cleaning HTML Forms (x-www-form-urlencoded), not JSON payloads

type HPPOptions struct {
	CheckQuery bool
	CheckParams bool
	Whitelist []string
	CheckBodyOnlyForContentType string
}

func Hpp(options HPPOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {

			if options.CheckParams && r.Method == http.MethodPost && isCorrectContentType(r, options.CheckBodyOnlyForContentType) {
				filterBodyParams(r, options.Whitelist)
			}

			if options.CheckQuery && r.URL.Query() != nil {
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
		fmt.Println(err)
		return	
	}

	for key, v := range r.Form {

		if len(v) > 1 {
			r.Form.Set(key, v[0])
		}
		if !isWhiteListed(key, whitelist) {
			delete(r.Form, key)
		}
	} 
}

func isWhiteListed(param string, whitelist []string) bool {
	for _, p := range whitelist {
		if p == param {
			return true
		}
	}
	return false
}

func filterQueryParams(r *http.Request, whitelist []string) {
	query := r.URL.Query()

	for key, v := range query {
		if len(v) > 1 {
			query.Set(key, v[0])
		}
		if !isWhiteListed(key, whitelist) {
			query.Del(key)
		}
	}

	r.URL.RawQuery = query.Encode() // Raw Query passed to next middleware or handler will have the cleaned query
}