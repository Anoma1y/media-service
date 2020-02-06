package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"./helpers"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	log.Println("Starting app")

	r := gin.Default()
	r.POST("/api/v0/file", uploadFile)
	r.POST("/api/v1/file", uploadFileS3)

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

	pathname := "http://localhost:3000" + "/" + groupName + "/" + filename // todo

	c.JSON(http.StatusCreated, gin.H{
		"filename":  filename,
		"extension": extension,
		"size":      size,
		"path":      pathname,
	})
}

func uploadFileS3(c *gin.Context) {
	accessKey := os.Getenv("AWS_SECRET_KEY")
	secretKey := os.Getenv("AWS_ACCESS_KEY")
	endpoint := os.Getenv("AWS_ENDPOINT")
	region := os.Getenv("AWS_REGION")
	bucket := os.Getenv("AWS_BUCKET")

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint),
		Credentials: credentials.NewStaticCredentials(
			secretKey,
			accessKey,
			"",
		),
	})

	if err != nil {
		fmt.Println("Could not upload file")
	}

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

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "[ERROR]: Saving file error" + "./Dockerfile"})
		return
	}

	extension, filename := helpers.GetFileMeta(header.Filename)

	groupName := helpers.GetGroupFileName()

	buffer := make([]byte, size)
	file.Read(buffer)
	fmt.Println(http.DetectContentType(buffer))
	_, s3err := s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(groupName + "/" + filename),
		ACL:           aws.String("public-read"),
		Body:          bytes.NewReader(buffer),
		ContentLength: aws.Int64(int64(size)),
		ContentType:   aws.String(http.DetectContentType(buffer)),
	})

	if s3err != nil {
		fmt.Println(s3err)
	}

	pathname := "https://" + bucket + "." + endpoint + "/" + groupName + "/" + filename // todo

	c.JSON(http.StatusCreated, gin.H{
		"filename":  filename,
		"extension": extension,
		"size":      size,
		"path":      pathname,
	})
}
