package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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

	filename := header.Filename
	groupName := getGroupFileName()

	dirname := fmt.Sprint("files/", groupName)

	newfilename := fmt.Sprint(dirname, "/"+generateFileName(24, getFileExt(filename)))

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

func getFileExt(filename string) string {
	fileArr := strings.Split(filename, ".")

	return fileArr[len(fileArr)-1]
}

func generateFileName(len int, ext string) string {
	randBytes := make([]byte, len)
	rand.Read(randBytes)

	return filepath.Join(hex.EncodeToString(randBytes) + "." + ext)
}
