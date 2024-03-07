package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func StartProxyServer(proxyPort string, targetPorts string) {
	targets := strings.Split(targetPorts, ",")
	proxies := make([]*httputil.ReverseProxy, len(targets))

	for i, port := range targets {
		targetURL, err := url.Parse("http://localhost:" + port)
		if err != nil {
			log.Fatal(err)
		}
		proxies[i] = httputil.NewSingleHostReverseProxy(targetURL)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		proxies[len(r.URL.Path)%len(proxies)].ServeHTTP(w, r)
	}

	server := &http.Server{
		Addr:    ":" + proxyPort,
		Handler: http.HandlerFunc(handler),
	}

	log.Printf("Starting proxy server on :%s\n", proxyPort)
	log.Fatal(server.ListenAndServe())
}
