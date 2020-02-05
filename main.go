package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"./helpers"
	"github.com/gin-gonic/gin"
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

	size := header.Size

	if size > (1024 * 1024 * 5) {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"status": "[ERROR] File is too big"})
		return
	}

	extension, filename := helpers.GetFileMeta(header.Filename)

	groupName := helpers.GetGroupFileName()

	dirname := fmt.Sprint("files/", groupName)

	newfilename := fmt.Sprint(dirname, "/"+filename)

	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		os.Mkdir(dirname, 0700)
	}

	out, err := os.Create(newfilename)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "[ERROR] Creating file error " + filename})
		return
	}

	defer out.Close()

	_, err = io.Copy(out, file)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "[ERROR]: Saving file error" + filename})
		return
	}

	pathname := "http://localhost:3000" + "/" + groupName + "/" + filename

	c.JSON(http.StatusCreated, gin.H{
		"filename":  filename,
		"extension": extension,
		"size":      size,
		"path":      pathname,
	})
}
