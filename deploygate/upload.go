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

func (r *Response) IsSuccess() bool {
	return r.SuccessResponse != nil
}

func Upload(token, userName, filePath, message, distributionKey, distributionName, releaseNote string, disableNotify *bool, visibility string) *Response {
	endPointUrl := fmt.Sprintf("https://deploygate.com/api/users/%s/apps", userName)

	body, contentType, err := buildRequestBody(filePath, message, distributionKey, distributionName, releaseNote, disableNotify, visibility)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", endPointUrl, body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	return parseResponse(res)
}

func buildRequestBody(filePath, message, distributionKey, distributionName, releaseNote string, disableNotify *bool, visibility string) (io.Reader, string, error) {
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)

	if err := appendFormFile(writer, "file", filePath); err != nil {
		return nil, "", err
	}

	if message != "" {
		if err := appendFormField(writer, "message", message); err != nil {
			return nil, "", err
		}
	}
	if distributionKey != "" {
		if err := appendFormField(writer, "distribution_key", distributionKey); err != nil {
			return nil, "", err
		}
	}
	if distributionName != "" {
		if err := appendFormField(writer, "distribution_name", distributionName); err != nil {
			return nil, "", err
		}
	}
	if releaseNote != "" {
		if err := appendFormField(writer, "release_note", releaseNote); err != nil {
			return nil, "", err
		}
	}
	if disableNotify != nil {
		if err := appendFormField(writer, "disableNotify", strconv.FormatBool(*disableNotify)); err != nil {
			return nil, "", err
		}
	}
	if visibility != "" {
		if err := appendFormField(writer, "visibility", visibility); err != nil {
			return nil, "", err
		}
	}

	writer.Close()

	return buffer, writer.FormDataContentType(), nil
}

func parseResponse(httpResponse *http.Response) *Response {
	decoder := json.NewDecoder(httpResponse.Body)
	response := &Response{}

	switch {
	case httpResponse.StatusCode == http.StatusOK:
		response.SuccessResponse = &SuccessResponse{}
		decoder.Decode(response.SuccessResponse)
	case 400 <= httpResponse.StatusCode && httpResponse.StatusCode < 500:
		fmt.Fprintln(os.Stderr, httpResponse.Status)
		response.ErrorResponse = &ErrorResponse{}
		decoder.Decode(response.ErrorResponse)
	default:
		fmt.Fprintln(os.Stderr, httpResponse.Status)
	}

	return response
}

func appendFormFile(mpart *multipart.Writer, fieldName, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileName := filepath.Base(filePath)
	fileWriter, err := mpart.CreateFormFile(fieldName, fileName)
	if err != nil {
		return err
	}

	if _, err := io.Copy(fileWriter, file); err != nil {
		return err
	}

	return nil
}

func appendFormField(mpart *multipart.Writer, fieldName string, content string) error {
	fieldWriter, err := mpart.CreateFormField(fieldName)
	if err != nil {
		return err
	}
	if _, err := fieldWriter.Write([]byte(content)); err != nil {
		return err
	}
	return nil
}
