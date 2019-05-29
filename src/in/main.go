package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/YuukiARIA/concourse-deploygate-resource/deploygate"
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

	metadataPath := filepath.Join(basePath, "metadata.json")
	file, err := os.Create(metadataPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, _ := json.Marshal(&app)
	if _, err := file.Write(data); err != nil {
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

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Fatal: source directory not given")
		os.Exit(1)
	}
	basePath := os.Args[1]

	var request GetRequest
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid request: %s\n", err)
		os.Exit(1)
	}

	response, err := PerformGet(request, basePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	json.NewEncoder(os.Stdout).Encode(response)
}
