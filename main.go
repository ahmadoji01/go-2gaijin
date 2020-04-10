package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kitalabs/go-2gaijin/router"
)

func homePage(c *gin.Context) {}

func main() {
	fmt.Println("Distributed Chat App v0.01")

	router := router.Router()

	log.Fatal(http.ListenAndServe(":8080", router))
}
