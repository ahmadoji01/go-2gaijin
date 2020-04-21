package main

import (
	"log"
	"net/http"

	"gitlab.com/kitalabs/go-2gaijin/router"
)

func main() {

	router := router.Router()
	//log.Fatal(http.ListenAndServe(":80", router))
	log.Fatal(http.ListenAndServeTLS(":443", "keys/cert.pem", "keys/key.pem", router))
}
