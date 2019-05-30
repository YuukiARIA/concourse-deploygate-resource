package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/YuukiARIA/concourse-deploygate-resource/deploygate"
	"github.com/YuukiARIA/concourse-deploygate-resource/logger"
	"github.com/YuukiARIA/concourse-deploygate-resource/models"
)

func readMessageFromFile(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func dg(request *models.PutRequest) (*deploygate.UploadResponse, error) {
	source, params := request.Source, request.Params

	message := params.Message
	if params.MessageFile != "" {
		var err error
		message, err = readMessageFromFile(params.MessageFile)
		if err != nil {
			return nil, err
		}
	}

	dgClient := deploygate.NewClient(source.ApiKey)
	dgResponse, err := dgClient.Upload(
		source.Owner,
		params.File,
		message,
		params.DistributionKey,
		params.DistributionName,
		params.ReleaseNote,
		params.DisableNotify,
		params.Visibility,
	)
	if err != nil {
		return nil, err
	}
	return dgResponse, nil
}

func main() {
	if len(os.Args) < 2 {
		logger.Fatal("Fatal: source directory not given")
	}

	os.Chdir(os.Args[1])

	request := models.PutRequest{}
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		logger.Fatalf("Fatal: %s\n", err)
	}

	dgResponse, err := dg(&request)
	if err != nil {
		logger.Fatalf("Fatal: %s\n", err)
	} else if dgResponse.Error {
		logger.Fatalf("Error: %s (%s)\n", dgResponse.Message, dgResponse.Because)
	}

	results := dgResponse.Results

	response := models.PutResponse{
		Version: models.Version{
			Platform:    results.OSName,
			PackageName: results.PackageName,
			Revision:    strconv.Itoa(results.Revision),
		},
		Metadata: []models.MetadataEntry{
			{Name: "name", Value: results.Name},
			{Name: "package", Value: results.PackageName},
			{Name: "platform", Value: results.OSName},
			{Name: "revision", Value: strconv.Itoa(results.Revision)},
			{Name: "version_code", Value: results.VersionCode},
			{Name: "version_name", Value: results.VersionName},
			{Name: "url", Value: "https://deploygate.com" + results.Path},
			{Name: "message", Value: results.Message},
		},
	}
	json.NewEncoder(os.Stdout).Encode(response)
}
