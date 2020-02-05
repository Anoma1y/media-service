package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"path/filepath"
	"strings"
	"time"
)

func GetFileMeta(Filename string) (string, string) {
	extension := GetFileExt(Filename)
	filename := GenerateFileName(18, extension)

	return extension, filename
}

func GetFileExt(filename string) string {
	fileArr := strings.Split(filename, ".")

	return fileArr[len(fileArr)-1]
}

func GenerateFileName(len int, ext string) string {
	randBytes := make([]byte, len)
	rand.Read(randBytes)

	return filepath.Join(hex.EncodeToString(randBytes) + "." + ext)
}

func GetGroupFileName() string {
	currentTime := time.Now()

	return currentTime.Format("20060102")
}
