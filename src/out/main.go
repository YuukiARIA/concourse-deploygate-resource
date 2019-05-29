package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/YuukiARIA/concourse-deploygate-resource/deploygate"
	"github.com/YuukiARIA/concourse-deploygate-resource/models"
)

func readMessageFromFile(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func dg(request *models.PutRequest) (*deploygate.Response, error) {
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
		fmt.Fprintln(os.Stderr, "Fatal: source directory not given")
		os.Exit(1)
	}

	os.Chdir(os.Args[1])

	request := models.PutRequest{}
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: %s\n", err)
		os.Exit(1)
	}

	dgResponse, err := dg(&request)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: %s\n", err)
		os.Exit(1)
	} else if dgResponse.Error {
		fmt.Fprintf(os.Stderr, "Error: %s (%s)\n", dgResponse.Message, dgResponse.Because)
		os.Exit(1)
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
