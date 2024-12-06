package main

const (
	EndPoint              = "192.168.86.6:9000"
	AccessKey             = "blcot6nCennqusvX7XYK"
	SecretKey             = "IpOiAiPapoUz4x1an873JlNG7XBq1Mfzpe2E0xK2"
	serviceName           = "FileStorageWithMinio"
	Ok                    = "OK"
	NotOk                 = "Not-Ok"
	MsgFileSaved          = "filed saved successfully"
	MsgFileDownloaded     = "file download successfully"
	MsgFileDeleted        = "file deleted successfully"
	ErrFileNotParsed      = "file Could not be parsed"
	ErrNoFileUploaded     = "Not file uploaded"
	ErrBucketNotFound     = "bucket not found"
	ErrFileNotFound       = "file not found"
	ErrDataNotFound       = "data not found"
	ErrDataCouldNotRead   = "data could not be read"
	ErrFileNotSaved       = "file could not be saved"
	ErrFileOpen           = "file could not be opened"
	ErrDecode             = "could not decode the json"
	ErrNotString          = "value is not a string"
	ErrNotEmptyString     = "value is empty"
	ErrMissingFieldBucket = "bucket field not passed"
	ErrMissingFieldName   = "fileName field not passed"
	ErrLocalFileCreate    = "could not create local file"
	ErrLocalFileWright    = "could not write data in local file"
)

var MsgCode = map[string]string{
	MsgFileSaved:          "INFO001",
	MsgFileDownloaded:     "INFO002",
	MsgFileDeleted:        "INFO003",
	ErrFileNotParsed:      "ERR001",
	ErrNoFileUploaded:     "ERR002",
	ErrFileNotSaved:       "ERR003",
	ErrFileOpen:           "ERR004",
	ErrDecode:             "ERR005",
	ErrNotString:          "ERR006",
	ErrNotEmptyString:     "ERR007",
	ErrMissingFieldBucket: "ERR008",
	ErrMissingFieldName:   "ERR009",
	ErrDataNotFound:       "ERR010",
	ErrDataCouldNotRead:   "ERR011",
	ErrLocalFileCreate:    "ERR012",
	ErrLocalFileWright:    "ERR013",
	ErrBucketNotFound:     "ERR014",
	ErrFileNotFound:       "ERR015",
}
