package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGiveResponse(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		response   Response
		want       int
	}{
		{
			name:       "200",
			statusCode: http.StatusOK,
			response: Response{
				ServiceName: serviceName,
			},
			want: http.StatusOK,
		},

		{
			name:       "400",
			statusCode: http.StatusBadRequest,
			response: Response{
				ServiceName: serviceName,
			},
			want: http.StatusBadRequest,
		},

		{
			name:       "500",
			statusCode: http.StatusInternalServerError,
			response: Response{
				ServiceName: serviceName,
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			giveResponse(w, tt.statusCode, tt.response)

			if w.Code != tt.want {
				t.Errorf("\ngot status code : %v\nwant status code : %v", w.Code, tt.want)
			}
		})
	}

}

func TestIdGenerator(t *testing.T) {
	testCases := []struct {
		name  string
		input int
		want  int
	}{
		{
			name:  "valid",
			input: 10,
			want:  10,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			id := IdGenerator(tt.input)
			gotId := []rune(id)
			if id == "" {
				t.Errorf("expected some string\ngot empty")
			}
			if len(gotId) != tt.want {
				t.Errorf("\ngot id of length %v\nwant id of lenght %v", len(gotId), tt.want)
			}
		})
	}

}

func TestFetchContentType(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  ".pdf",
			input: ".pdf",
			want:  "application/pdf",
		},
		{
			name:  ".jpeg",
			input: ".jpeg",
			want:  "image/jpg",
		},
		{
			name:  ".jpg",
			input: ".jpg",
			want:  "image/jpg",
		},
		{
			name:  ".png",
			input: ".png",
			want:  "image/png",
		},
		{
			name:  ".json",
			input: ".json",
			want:  "application/json",
		},
		{
			name:  ".docx",
			input: ".docx",
			want:  "application/docx",
		},
		{
			name:  ".docs",
			input: ".docs",
			want:  "application/docx",
		},
		{
			name:  ".doc",
			input: ".doc",
			want:  "application/docx",
		},
		{
			name:  "unknown",
			input: ".exe",
			want:  "",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := fetchContentType(tt.input)
			if got != tt.want {
				t.Errorf("\ngot : %v\n want : %v", got, tt.want)
			}
		})
	}
}

func TestReqValidation(t *testing.T) {
	testCases := []struct {
		name  string
		input map[string]interface{}
		want  error
	}{
		{
			name: "valid",
			input: map[string]interface{}{
				"bucket":   "practice",
				"fileName": "hello.pdf",
			},
			want: nil,
		},
		{
			name: "bucket field missing",
			input: map[string]interface{}{
				"hello":    "practice",
				"fileName": "hello.pdf",
			},
			want: errors.New(ErrMissingFieldBucket),
		},
		{
			name: "bucket not string",
			input: map[string]interface{}{
				"bucket":   123456,
				"fileName": "hello.pdf",
			},
			want: errors.New(ErrNotString),
		},
		{
			name: "bucket empty",
			input: map[string]interface{}{
				"bucket":   "",
				"fileName": "hello.pdf",
			},
			want: errors.New(ErrNotEmptyString),
		},
		{
			name: "fileName field missing",
			input: map[string]interface{}{
				"bucket": "practice",
				"hello":  "hello.pdf",
			},
			want: errors.New(ErrMissingFieldName),
		},
		{
			name: "fileName not string",
			input: map[string]interface{}{
				"bucket":   "practice",
				"fileName": 123456,
			},
			want: errors.New(ErrNotString),
		},
		{
			name: "fileName empty",
			input: map[string]interface{}{
				"bucket":   "practice",
				"fileName": "",
			},
			want: errors.New(ErrNotEmptyString),
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := reqValidation(tt.input)
			if got != nil {
				if got.Error() != tt.want.Error() {
					t.Errorf("\ngot error : %v\nwant error %v", got.Error(), tt.want.Error())
				}
			}
		})
	}
}
