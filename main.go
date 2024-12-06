package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	mx := http.NewServeMux()

	mx.Handle("POST /FileStorageWithMinio/upload", http.HandlerFunc(UploadFile))
	mx.Handle("POST /FileStorageWithMinio/download", http.HandlerFunc(DownloadFile))
	mx.Handle("POST /FileStorageWithMinio/delete", http.HandlerFunc(DeleteFile))

	fmt.Println("Server started on port 1333")
	log.Fatal(http.ListenAndServe(":1333", mx))
}
