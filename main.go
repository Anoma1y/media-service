package main

import (
  "log"
  "github.com/gin-gonic/gin"
)

const (
  StatusOK = 200
  StatusCreated = 201
  
  StatusBadRequest = 400
  StatusForbidden = 403
  StatusNotFound = 404
  
  StatusInternalServerError = 500
)

func main() {
  log.Println("Starting app")
  
  r := gin.Default()
  
  // r.POST("/api/v1/upload", uploadFile)
  // r.POST("/api/v1/image/resize", resizeImage)
  // r.GET("/api/v1/image", getImage)
  
  port := "3000"
  
  r.Run(":" + port)
  
  log.Println("App is ready")
}

// func uploadFile(c *gin.Context) {

// }

// func resizeImage(c *gin.Context) {

// }

// func getImage(c *gin.Context) {

// }
      