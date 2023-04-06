package middlewares

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r != nil {
			rawBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				log.Printf("Request: %v, %v, %v", r.Method, r.RequestURI, string(rawBody))
			} else {
				log.Printf("Request: %v, %v, %v", r.Method, r.RequestURI, err.Error())
			}
			r.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))
		}
		next.ServeHTTP(w, r)
	})
}
