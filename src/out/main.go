package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

func dg(request *Request) *deploygate.Response {
	source, params := request.Source, request.Params

	message := params.Message
	if params.MessageFile != "" {
		var err error
		message, err = readMessageFromFile(params.MessageFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	return deploygate.Upload(
		source.ApiKey,
		source.Owner,
		params.File,
		message,
		params.DistributionKey,
		params.DistributionName,
		params.ReleaseNote,
		params.DisableNotify,
		params.Visibility,
	)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "error: source directory not given")
		os.Exit(1)
	}

	os.Chdir(os.Args[1])

	request := Request{}
	json.NewDecoder(os.Stdin).Decode(&request)

	dgResponse := dg(&request)
	if !dgResponse.IsSuccess() {
		fmt.Fprintf(os.Stderr, "error message=%s, because=%s\n", dgResponse.ErrorResponse.Message, dgResponse.ErrorResponse.Because)
		os.Exit(1)
	}

	results := dgResponse.SuccessResponse.Results

	response := Response{
		Version: Version{},
		Metadata: []MetadataEntry{
			MetadataEntry{Name: "name", Value: results.Name},
			MetadataEntry{Name: "package", Value: results.PackageName},
			MetadataEntry{Name: "platform", Value: results.OSName},
			MetadataEntry{Name: "revision", Value: strconv.Itoa(results.Revision)},
			MetadataEntry{Name: "version_code", Value: strconv.Itoa(results.VersionCode)},
			MetadataEntry{Name: "version_name", Value: results.VersionName},
			MetadataEntry{Name: "url", Value: "https://deploygate.com" + results.Path},
			MetadataEntry{Name: "message", Value: results.Message},
		},
	}
	json.NewEncoder(os.Stdout).Encode(response)
}
