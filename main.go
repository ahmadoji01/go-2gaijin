package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kitalabs/go-2gaijin/router"
)

func homePage(c *gin.Context) {}

func main() {

	router := router.Router()

	log.Fatal(http.ListenAndServeTLS(":8080", "keys/cert.pem", "keys/key.pem", router))
}
