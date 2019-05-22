package main

import (
	"encoding/json"
	"fmt"
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
	File        string `json:"file"`
	Message     string `json:"message"`
	MessageFile string `json:"message_file"`
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

func dg(request *Request) *deploygate.Response {
	return deploygate.Upload(
		request.Source.ApiKey,
		request.Source.Owner,
		request.Params.File,
		request.Params.Message,
		"",
		"",
		"",
		nil,
		"",
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

	success := dgResponse.SuccessResponse

	response := Response{
		Version: Version{},
		Metadata: []MetadataEntry{
			MetadataEntry{Name: "revision", Value: strconv.Itoa(success.Results.Revision)},
			MetadataEntry{Name: "url", Value: "https://deploygate.com" + success.Results.Path},
		},
	}
	json.NewEncoder(os.Stdout).Encode(response)
}
