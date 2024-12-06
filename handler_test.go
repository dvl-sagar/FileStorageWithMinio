package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/minio/minio-go/v7"
)

func TestUploadFile(t *testing.T) {
	testCases := []struct {
		name     string
		files    []string
		wantCode int
	}{
		{
			name:     "Single file upload",
			files:    []string{"tempDocxFile.docx"},
			wantCode: http.StatusOK,
		},
		{
			name:     "Multiple files upload",
			files:    []string{"tempDocxFile.docx", "tempfiletwo.docx"},
			wantCode: http.StatusOK,
		},
		{
			name:     "no files uploaded",
			files:    []string{},
			wantCode: http.StatusBadRequest,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Create a multipart form with the specified files
			var requestBody bytes.Buffer
			writer := multipart.NewWriter(&requestBody)

			for _, filePath := range tt.files {
				file, _ := os.Open(filePath)
				defer file.Close()

				part, err := writer.CreateFormFile("file", filePath)
				if err != nil {
					t.Fatalf("Error creating form file part: %v", err)
				}
				_, _ = io.Copy(part, file)
				if err != nil {
					t.Fatalf("Error copying file content: %v", err)
				}
			}
			// Close the writer
			err := writer.Close()
			if err != nil {
				t.Fatalf("Error closing writer: %v", err)
			}

			r := httptest.NewRequest("POST", "/upload", &requestBody)
			r.Header.Set("Content-Type", writer.FormDataContentType())

			w := httptest.NewRecorder()
			UploadFile(w, r)
			if w.Code != tt.wantCode {
				t.Errorf("\ngot code : %v\nwant code : %v", w.Code, tt.wantCode)
			}
		})
		t.Run("invalid value", func(t *testing.T) {
			// Create a multipart form with the specified files
			var requestBody bytes.Buffer
			writer := multipart.NewWriter(&requestBody)

			file, _ := os.Open("TempFile.pdf")
			defer file.Close()

			part, err := writer.CreateFormFile("file", "TempFile.pdf")
			if err != nil {
				t.Fatalf("Error creating form file part: %v", err)
			}
			_, _ = io.Copy(part, file)
			if err != nil {
				t.Fatalf("Error copying file content: %v", err)
			}

			// Close the writer
			err = writer.Close()
			if err != nil {
				t.Fatalf("Error closing writer: %v", err)
			}

			r := httptest.NewRequest("POST", "/upload", &requestBody)
			r.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			UploadFile(w, r)
			if w.Code != http.StatusBadRequest {
				t.Errorf("\ngot code : %v\nwant code : %v", w.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestDownloadFile(t *testing.T) {
	testCases := []struct {
		name     string
		data     interface{}
		wantCode int
	}{
		{
			name:     "invalid JSON",
			data:     `{"name": "John"`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "validation error bucket field not passed",
			data: map[string]interface{}{
				"name":     "practice",
				"fileName": "something.pdf",
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "validation error fileName field not passed",
			data: map[string]interface{}{
				"bucket": "practice",
				"name":   "something.pdf",
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "file not present",
			data: map[string]interface{}{
				"bucket":   "practice",
				"fileName": "something.pdf",
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "bucket not present",
			data: map[string]interface{}{
				"bucket":   "practices",
				"fileName": "something.pdf",
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "valid",
			data: map[string]interface{}{
				"bucket":   "practice",
				"fileName": "TempFile_QIKBIBDGVGHF.pdf",
			},
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			inputJson, err := json.Marshal(tt.data)
			if err != nil {
				t.Errorf("failed to marshal input : %v", err)
			}

			r := httptest.NewRequest("POST", "/FileStorageWithMinio/download", bytes.NewBuffer(inputJson))
			w := httptest.NewRecorder()

			DownloadFile(w, r)
			if w.Code != tt.wantCode {
				t.Errorf("\ngot status code : %v\nwant status code : %v", w.Code, tt.wantCode)
			}
		})
	}
}

func TestDeleteFile(t *testing.T) {
	connectMinio()
	fileData, err := os.Open("demofile.pdf")
	if err != nil {
		t.Errorf("err reading demofile")
	}
	bucketName := "practice"
	fileId := IdGenerator(12)

	ext := filepath.Ext(fileData.Name())
	fileName := strings.TrimSuffix(fileData.Name(), ext)
	FinalFileName := fileName + "_" + fileId + ext
	contentType := fetchContentType(ext)
	fileObject, _ := fileData.Stat()
	fmt.Println(contentType)

	_, _ = minioClient.PutObject(context.Background(), bucketName, FinalFileName, fileData, fileObject.Size(), minio.PutObjectOptions{ContentType: contentType})

	testCases := []struct {
		name     string
		data     interface{}
		wantCode int
	}{
		{
			name:     "invalid JSON",
			data:     `{"name": "John"`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "validation error bucket field not passed",
			data: map[string]interface{}{
				"name":     "practice",
				"fileName": "something.pdf",
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "validation error fileName field not passed",
			data: map[string]interface{}{
				"bucket": "practice",
				"name":   "something.pdf",
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "file not present",
			data: map[string]interface{}{
				"bucket":   "practice",
				"fileName": "something.pdf",
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "bucket not present",
			data: map[string]interface{}{
				"bucket":   "practices",
				"fileName": "something.pdf",
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "valid",
			data: map[string]interface{}{
				"bucket":   "practice",
				"fileName": FinalFileName,
			},
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			inputJson, err := json.Marshal(tt.data)
			if err != nil {
				t.Errorf("failed to marshal input : %v", err)
			}
			r := httptest.NewRequest("POST", "/delete", bytes.NewBuffer(inputJson))
			w := httptest.NewRecorder()

			DeleteFile(w, r)
			if w.Code != tt.wantCode {
				t.Errorf("\ngot status code : %v\nwant status code : %v", w.Code, tt.wantCode)
			}
		})
	}
}
