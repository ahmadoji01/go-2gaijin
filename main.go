package main

import (
	"log"
	"net/http"

	"gitlab.com/kitalabs/go-2gaijin/router"
)

func main() {

	router := router.Router()
	log.Fatal(http.ListenAndServe(":8080", router))
}
