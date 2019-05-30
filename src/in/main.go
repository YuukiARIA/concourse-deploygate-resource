package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"

	"github.com/YuukiARIA/concourse-deploygate-resource/deploygate"
	"github.com/YuukiARIA/concourse-deploygate-resource/logger"
	"github.com/YuukiARIA/concourse-deploygate-resource/models"
)

type GetRequest struct {
	Source  *models.Source  `json:"source"`
	Version *models.Version `json:"version"`
	Params  *GetParams      `json:"params"`
}

type GetResponse struct {
	Version models.Version `json:"version"`
}

type GetParams struct {
}

func PerformGet(request GetRequest, basePath string) (*GetResponse, error) {
	dgClient := deploygate.NewClient(request.Source.ApiKey)
	getAppRes, err := dgClient.GetApp(request.Source.Owner, request.Version.Platform, request.Version.PackageName)
	if err != nil {
		return nil, err
	}

	app := getAppRes.App

	data, err := json.Marshal(&app)
	if err != nil {
		return nil, err
	}
	if err := writeFile(basePath, "metadata.json", data); err != nil {
		return nil, err
	}

	return &GetResponse{
		Version: models.Version{
			Platform:    app.OSName,
			PackageName: app.PackageName,
			Revision:    strconv.Itoa(app.CurrentRevision),
		},
	}, nil
}

func writeFile(basePath, fileName string, content []byte) error {
	filePath := filepath.Join(basePath, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(content); err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		logger.Fatal("Fatal: source directory not given")
	}
	basePath := os.Args[1]

	var request GetRequest
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		logger.Fatalf("Invalid request: %s\n", err)
	}

	response, err := PerformGet(request, basePath)
	if err != nil {
		logger.Fatal(err)
	}
	json.NewEncoder(os.Stdout).Encode(response)
}
