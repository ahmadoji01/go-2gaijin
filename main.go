package main

import (
	"log"
	"net/http"

	"gitlab.com/kitalabs/go-2gaijin/router"
)

func main() {

	router := router.Router()
	err := http.ListenAndServeTLS(":443", "server.crt", "server.key", router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
