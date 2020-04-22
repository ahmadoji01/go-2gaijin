package main

import (
	"log"
	"net/http"

	"gitlab.com/kitalabs/go-2gaijin/router"
)

const (
	domainName = "go.2gaijin.com"
)

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+domainName+":443"+r.RequestURI, http.StatusMovedPermanently)
}

func main() {

	router := router.Router()
	go func() {
		if err := http.ListenAndServe(":80", http.HandlerFunc(redirectTLS)); err != nil {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()
	log.Fatal(http.ListenAndServeTLS(":443", "keys/cert.pem", "keys/key.pem", router))
}
