package deploygate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func (c *Client) Upload(userName, filePath, message, distributionKey, distributionName, releaseNote string, disableNotify bool, visibility string) (*UploadResponse, error) {
	endPointUrl := fmt.Sprintf("https://deploygate.com/api/users/%s/apps", userName)

	body, contentType, err := buildRequestBody(filePath, message, distributionKey, distributionName, releaseNote, disableNotify, visibility)
	if err != nil {
		return nil, err
	}

	req, err := c.buildRequest(endPointUrl, body, contentType)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return parseResponse(res)
}

func (c *Client) buildRequest(url string, body io.Reader, contentType string) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.ApiKey)
	req.Header.Set("Content-Type", contentType)
	return req, nil
}

func buildRequestBody(filePath, message, distributionKey, distributionName, releaseNote string, disableNotify bool, visibility string) (io.Reader, string, error) {
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
	if disableNotify {
		if err := appendFormField(writer, "disableNotify", "yes"); err != nil {
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

func parseResponse(httpResponse *http.Response) (*UploadResponse, error) {
	response := &UploadResponse{}
	if err := json.NewDecoder(httpResponse.Body).Decode(response); err != nil {
		return nil, err
	}
	return response, nil
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
