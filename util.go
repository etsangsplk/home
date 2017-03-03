package main

import (
	"net/http"
	"net/url"
	"strings"
)

// copyRequestAndURL returns a copy of req and its URL field.
func copyRequestAndURL(req *http.Request) *http.Request {
	r := *req
	u := *req.URL
	r.URL = &u
	return &r
}

// http۰StripPrefix is http.StripPrefix from Go 1.9, with https://github.com/golang/go/issues/18952 fixed.
// TODO: Replace with http.StripPrefix once Go 1.9 comes out.
func http۰StripPrefix(prefix string, h http.Handler) http.Handler {
	if prefix == "" {
		return h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if p := strings.TrimPrefix(r.URL.Path, prefix); len(p) < len(r.URL.Path) {
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = p
			h.ServeHTTP(w, r2)
		} else {
			http.NotFound(w, r)
		}
	})
}