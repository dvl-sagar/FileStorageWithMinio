package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

func connectMinio() {
	var err error

	minioClient, err = minio.New(EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(AccessKey, SecretKey, ""),
		Secure: false})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("###-Minio Connected Successfully-###")
}

func giveResponse(w http.ResponseWriter, statusCode int, response Response) {
	w.WriteHeader(statusCode)
	res, _ := json.Marshal(response)
	w.Write(res)
}

func IdGenerator(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.NewSource(time.Now().UnixNano())
	b := make([]rune, n)
	fmt.Println("before : " + string(b))
	for i := range b {
		b[i] = letter[rand.Intn(int(len(letter)))]
	}
	fmt.Println("After : " + string(b))
	return strings.ToUpper(string(b))
}

func fetchContentType(ext string) string {
	switch ext {
	case ".pdf":
		return "application/pdf"
	case ".jpeg", ".jpg":
		return "image/jpg"
	case ".png":
		return "image/png"
	case ".json":
		return "application/json"
	case ".docx", ".docs", ".doc":
		return "application/docx"
	default:
		return ""
	}
}

func reqValidation(req map[string]interface{}) error {
	if bucket, ok := req["bucket"]; ok {
		if str, ok := bucket.(string); ok {
			if str == "" {
				return errors.New(ErrNotEmptyString)
			}
		} else {
			return errors.New(ErrNotString)
		}
	} else {
		return errors.New(ErrMissingFieldBucket)
	}
	if name, ok := req["fileName"]; ok {
		if str, ok := name.(string); ok {
			if str == "" {
				return errors.New(ErrNotEmptyString)
			}
		} else {
			return errors.New(ErrNotString)
		}
	} else {
		return errors.New(ErrMissingFieldName)
	}
	return nil
}
