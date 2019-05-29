package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/YuukiARIA/concourse-deploygate-resource/deploygate"
	"github.com/YuukiARIA/concourse-deploygate-resource/models"
)

type CheckRequest struct {
	Source  models.Source  `json:"source"`
	Version models.Version `json:"version"`
}

type CheckResponse []models.Version

func PerformCheck(request CheckRequest) (*CheckResponse, error) {
	dgClient := deploygate.NewClient(request.Source.ApiKey)
	getAppsRes, err := dgClient.GetApps(request.Source.Owner)
	if err != nil {
		return nil, err
	}

	versions := CheckResponse{}
	for _, app := range getAppsRes.Apps {
		version := models.Version{
			Platform:    app.OSName,
			PackageName: app.PackageName,
			Revision:    strconv.Itoa(app.CurrentRevision),
		}
		versions = append(versions, version)
	}

	return &versions, nil
}

func main() {
	var request CheckRequest
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	response, err := PerformCheck(request)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	json.NewEncoder(os.Stdout).Encode(response)
}
