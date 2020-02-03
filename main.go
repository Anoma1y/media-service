package main

import (
	// "fmt"
	"log"
	"github.com/gin-gonic/gin"
	// "os"
)

func main() {
	log.Println("Starting app....")

	r := gin.Default()

	port := "3000"

	r.Run(":" + port)
}
