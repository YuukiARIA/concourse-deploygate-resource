package main

import (
	"encoding/json"
	"fmt"
	"os"

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

func dg(request *Request) {
	deploygate.Upload(
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

	dg(&request)

	response := Response{
		Version: Version{},
		Metadata: []MetadataEntry{
			MetadataEntry{Name: "message", Value: request.Params.Message},
		},
	}
	json.NewEncoder(os.Stdout).Encode(response)
}
