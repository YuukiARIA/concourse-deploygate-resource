package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/YuukiARIA/concourse-deploygate-resource/deploygate"
)

type Version struct {
}

type Source struct {
	ApiKey string `json:"api_key"`
	Owner  string `json:"owner"`
}

type Params struct {
	File             string `json:"file"`
	Message          string `json:"message"`
	MessageFile      string `json:"message_file"`
	ReleaseNote      string `json:"release_note"`
	DistributionKey  string `json:"distribution_key"`
	DistributionName string `json:"distribution_name"`
	DisableNotify    bool   `json:"disable_notify"`
	Visibility       string `json:"visibility"`
}

type Request struct {
	Source Source `json:"source"`
	Params Params `json:"params"`
}

type MetadataEntry struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Response struct {
	Version  Version         `json:"version"`
	Metadata []MetadataEntry `json:"metadata"`
}

func readMessageFromFile(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func dg(request *Request) (*deploygate.Response, error) {
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

	request := Request{}
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

	response := Response{
		Version: Version{},
		Metadata: []MetadataEntry{
			{Name: "name", Value: results.Name},
			{Name: "package", Value: results.PackageName},
			{Name: "platform", Value: results.OSName},
			{Name: "revision", Value: strconv.Itoa(results.Revision)},
			{Name: "version_code", Value: strconv.Itoa(results.VersionCode)},
			{Name: "version_name", Value: results.VersionName},
			{Name: "url", Value: "https://deploygate.com" + results.Path},
			{Name: "message", Value: results.Message},
		},
	}
	json.NewEncoder(os.Stdout).Encode(response)
}
