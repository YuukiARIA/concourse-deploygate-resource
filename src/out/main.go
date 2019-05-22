package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Version struct {
}

type Request struct {
}

type Response struct {
	Version Version `json:"version"`
}

func main() {
	fmt.Fprintln(os.Stderr, "(do nothing)")
	json.NewEncoder(os.Stdout).Encode(Response{Version: Version{}})
}
