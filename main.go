package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/elazarl/goproxy"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	// Configure the proxy to intercept HTTPS traffic
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	proxy.OnRequest().DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		path := r.URL.Path
		log.Printf("Intercepted path: %s", path)

		intercept, err := shouldIntercept(path)
		if err != nil {
			log.Printf("Error evaluating expression: %v", err)
			return r, nil
		}

		if intercept {
			customResponse := []byte("Custom JSON response")
			return r, goproxy.NewResponse(r,
				goproxy.ContentTypeText, http.StatusOK, string(customResponse))
		}

		return r, nil
	})

	log.Fatal(http.ListenAndServe(":4444", proxy))
}

func shouldIntercept(path string) (bool, error) {
	log.Printf("Intercepted path: %s", path)

	if strings.Contains(path, "reddit") {
		return true, nil
	}

	return false, nil
}
