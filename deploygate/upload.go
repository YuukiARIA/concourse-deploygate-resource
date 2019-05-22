package deploygate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Results struct {
	Name        string `json:"name"`
	PackageName string `json:"package_name"`
	OSName      string `json:"os_name"`
	Path        string `json:"path"`
	Revision    int    `json:"revision"`
	VersionCode int    `json:"version_code"`
	VersionName string `json:"version_name"`
	Message     string `json:"message"`
}

type ErrorResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Because string `json:"because"`
}

type SuccessResponse struct {
	Error   bool    `json:"error"`
	Results Results `json:"results"`
}

type Response struct {
	SuccessResponse *SuccessResponse
	ErrorResponse   *ErrorResponse
}

func (r *Response) IsSuccess() bool {
	return r.SuccessResponse != nil
}

func Upload(token, userName, filePath, message, distributionKey, distributionName, releaseNote string, disableNotify *bool, visibility string) *Response {
	endPointUrl := fmt.Sprintf("https://deploygate.com/api/users/%s/apps", userName)

	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)

	appendFormFile(writer, "file", filePath)

	if message != "" {
		appendFormField(writer, "message", message)
	}
	if distributionKey != "" {
		appendFormField(writer, "distribution_key", distributionKey)
	}
	if distributionName != "" {
		appendFormField(writer, "distribution_name", distributionName)
	}
	if releaseNote != "" {
		appendFormField(writer, "release_note", releaseNote)
	}
	if disableNotify != nil {
		appendFormField(writer, "disableNotify", strconv.FormatBool(*disableNotify))
	}
	if visibility != "" {
		appendFormField(writer, "visibility", visibility)
	}

	writer.Close()

	req, _ := http.NewRequest("POST", endPointUrl, buffer)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	return parseResponse(res)
}

func parseResponse(httpResponse *http.Response) *Response {
	decoder := json.NewDecoder(httpResponse.Body)
	response := &Response{}

	switch httpResponse.StatusCode {
	case http.StatusOK:
		response.SuccessResponse = &SuccessResponse{}
		decoder.Decode(response.SuccessResponse)
	case http.StatusBadRequest:
		response.ErrorResponse = &ErrorResponse{}
		decoder.Decode(response.ErrorResponse)
	default:
		fmt.Fprintf(os.Stderr, "unsupported status: %s\n", httpResponse.Status)
	}

	return response
}

func appendFormFile(mpart *multipart.Writer, fieldName, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileName := filepath.Base(filePath)
	fileWriter, err := mpart.CreateFormFile(fieldName, fileName)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := io.Copy(fileWriter, file); err != nil {
		log.Fatal(err)
	}
}

func appendFormField(mpart *multipart.Writer, fieldName string, content string) {
	fieldWriter, err := mpart.CreateFormField(fieldName)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := fieldWriter.Write([]byte(content)); err != nil {
		log.Fatal(err)
	}
}