package main

type Response struct {
	ServiceName string      `json:"serviceName,omitempty"`
	MessageCode string      `json:"messageCode,omitempty"`
	Status      string      `json:"status,omitempty"`
	Msg         string      `json:"msg,omitempty"`
	Data        interface{} `json:"data,omitempty"`
}

type FileResp struct {
	FileName         string `json:"fileName,omitempty"`
	OriginalFileName string `json:"originalFileName,omitempty"`
	MessageCode      string `json:"messageCode,omitempty"`
	Location         string `json:"location,omitempty"`
	Err              string `json:"err,omitempty"`
}
