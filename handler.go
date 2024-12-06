package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		res := Response{
			ServiceName: serviceName,
			MessageCode: MsgCode[ErrFileNotParsed],
			Status:      NotOk,
			Msg:         err.Error(),
		}
		giveResponse(w, http.StatusBadRequest, res)
		return
	}
	files := r.MultipartForm.File["file"]
	if len(files) == 0 {
		res := Response{
			ServiceName: serviceName,
			MessageCode: MsgCode[ErrNoFileUploaded],
			Status:      NotOk,
			Msg:         ErrNoFileUploaded,
		}
		giveResponse(w, http.StatusBadRequest, res)
	}
	var res []FileResp

	for _, file := range files {
		originalFileName := file.Filename
		fileData, _ := os.Open(file.Filename)

		connectMinio()

		bucketName := "practice"
		fileId := IdGenerator(12)

		ext := filepath.Ext(file.Filename)
		fileName := strings.TrimSuffix(file.Filename, ext)
		file.Filename = fileName + "_" + fileId + ext

		contentType := fetchContentType(ext)

		fmt.Println(contentType)
		result, err := minioClient.PutObject(context.Background(), bucketName, file.Filename, fileData, file.Size, minio.PutObjectOptions{ContentType: contentType})
		if err != nil {
			res = append(res, FileResp{
				FileName:         file.Filename,
				OriginalFileName: originalFileName,
				MessageCode:      MsgCode[ErrFileNotSaved],
				Err:              err.Error(),
			})
		}
		res = append(res, FileResp{
			FileName:         file.Filename,
			OriginalFileName: originalFileName,
			Location:         result.Bucket,
			MessageCode:      MsgCode[MsgFileSaved],
		})
	}
	resp := Response{
		ServiceName: serviceName,
		MessageCode: MsgCode[MsgFileSaved],
		Status:      Ok,
		Msg:         MsgFileSaved,
		Data:        res,
	}
	giveResponse(w, http.StatusOK, resp)
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		res := Response{
			ServiceName: serviceName,
			MessageCode: MsgCode[ErrDecode],
			Status:      NotOk,
			Msg:         err.Error(),
		}
		giveResponse(w, http.StatusBadRequest, res)
		return
	}
	err = reqValidation(req)
	if err != nil {
		res := Response{
			ServiceName: serviceName,
			MessageCode: MsgCode[err.Error()],
			Status:      NotOk,
			Msg:         err.Error(),
		}
		giveResponse(w, http.StatusBadRequest, res)
		return
	}
	bucket := req["bucket"].(string)
	fileName := req["fileName"].(string)
	fmt.Printf("\nbucketName : %v\nfileName : %v", bucket, fileName)

	connectMinio()

	obj, _ := minioClient.GetObject(context.Background(), bucket, fileName, minio.GetObjectOptions{})

	data, err := io.ReadAll(obj)
	if err != nil {
		res := Response{
			ServiceName: serviceName,
			MessageCode: MsgCode[ErrDataNotFound],
			Status:      NotOk,
			Msg:         err.Error(),
		}
		giveResponse(w, http.StatusInternalServerError, res)
		return
	}
	local, _ := os.Create("Download_" + fileName)
	defer local.Close()
	local.Write(data)

	res := Response{
		ServiceName: serviceName,
		MessageCode: MsgCode[MsgFileDownloaded],
		Status:      NotOk,
		Msg:         MsgFileDownloaded,
		Data:        local,
	}
	giveResponse(w, http.StatusOK, res)
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		res := Response{
			ServiceName: serviceName,
			MessageCode: MsgCode[ErrDecode],
			Status:      NotOk,
			Msg:         err.Error(),
		}
		giveResponse(w, http.StatusBadRequest, res)
		return
	}
	err = reqValidation(req)
	if err != nil {
		res := Response{
			ServiceName: serviceName,
			MessageCode: MsgCode[err.Error()],
			Status:      NotOk,
			Msg:         err.Error(),
		}
		giveResponse(w, http.StatusBadRequest, res)
		return
	}
	bucket := req["bucket"].(string)
	fileName := req["fileName"].(string)
	fmt.Printf("\nbucketName : %v\nfileName : %v", bucket, fileName)

	connectMinio()
	exist, _ := minioClient.BucketExists(context.Background(), bucket)
	if !exist {
		res := Response{
			ServiceName: serviceName,
			MessageCode: MsgCode[ErrBucketNotFound],
			Status:      NotOk,
			Msg:         ErrBucketNotFound,
		}
		giveResponse(w, http.StatusInternalServerError, res)
		return
	}
	_, err = minioClient.StatObject(context.Background(), bucket, fileName, minio.GetObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			res := Response{
				ServiceName: serviceName,
				MessageCode: MsgCode[ErrFileNotFound],
				Status:      NotOk,
				Msg:         err.Error(),
			}
			giveResponse(w, http.StatusInternalServerError, res)
			return
		}
	}

	err = minioClient.RemoveObject(context.Background(), bucket, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		res := Response{
			ServiceName: serviceName,
			MessageCode: MsgCode[ErrDataNotFound],
			Status:      NotOk,
			Msg:         err.Error(),
		}
		giveResponse(w, http.StatusInternalServerError, res)
		return
	}
	res := Response{
		ServiceName: serviceName,
		MessageCode: MsgCode[MsgFileDeleted],
		Status:      Ok,
		Msg:         MsgFileDeleted,
	}
	giveResponse(w, http.StatusOK, res)
}
