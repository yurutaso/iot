package iot

import (
	"crypto/tls"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
)

const (
	CACHE_DIR string = `/usr/local/etc/webhook-certs`
)

func HttpsServer(domain string, handler func(http.ResponseWriter, *http.Request)) {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain),
		Cache:      autocert.DirCache(CACHE_DIR),
	}

	http.HandleFunc("/", handler)

	server := &http.Server{
		Addr: ":https",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}
	go http.ListenAndServe(":http", certManager.HTTPHandler(nil))
	log.Fatal(server.ListenAndServeTLS("", ""))
}
