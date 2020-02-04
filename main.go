package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	log.Println("Starting app")

	r := gin.Default()

	r.POST("/api/v1/file", uploadFile)

	port := "3000"

	r.Run(":" + port)

	log.Println("App is ready")
}

func uploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "File not informed"})
		return
	}

	filename := header.Filename
	groupName := getGroupFileName()

	dirname := fmt.Sprint("files/", groupName)
	fmt.Println(getFileUUIDName(filename))
	newfilename := fmt.Sprint(dirname, "/"+getFileUUIDName(filename))

	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		os.Mkdir(dirname, 0700)
	}

	out, err := os.Create(newfilename)
	if err != nil {
		log.Println("[ERROR] Creating file error", err)
		c.JSON(http.StatusNotFound, gin.H{"status": "[ERROR] Creating file error " + filename})
		return
	}

	defer out.Close()

	_, err = io.Copy(out, file)

	if err != nil {
		log.Println("[ERROR]: Saving file error", err)
		c.JSON(http.StatusNotFound, gin.H{"status": "[ERROR]: Saving file error" + filename})
		return
	}
}

func getGroupFileName() string {
	currentTime := time.Now()

	return currentTime.Format("20060102")
}

func getFileUUIDName(filename string) string {
	asd := strings.Split(filename, ".")
	ext := asd[len(asd)-1]

	name := uuid.New().String() + "." + ext

	return name
}
