package main

import (
	"encoding/json"
	"fmt"
	resource "github.com/lorands/maven-get-resource"
	"github.com/lorands/maven-get-resource/in"
	"log"
	"os"
	"path/filepath"
)


func main() {
	var request in.Request

	destinationDir := os.Args[1]
	resourceDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		_ = fmt.Errorf("fail to get resource directory! Bailing out: %v", err)
		os.Exit(2)
	}

	inputRequest(&request)

	//validate
	if len(request.Version.Version) <= 0 {
		_ = fmt.Errorf("missing version")
		os.Exit(3)
	}

	if len(request.Source.Artifact) <= 0 {
		_ = fmt.Errorf("missing artifact definition")
		os.Exit(4)
	}

	if len(request.Source.URL) <= 0 {
		_ = fmt.Errorf("missing repository URL")
		os.Exit(5)
	}

	version := request.Version.Version

	inResource := resource.InResource{
		Source:         request.Source,
		Version:        request.Version,
		DestinationDir: destinationDir,
		ResourceDir:    resourceDir,
	}

	if err := inResource.Download(); err != nil {
		fatal("Fail to process.", err)
	}

	response := in.Response{
		Version: resource.Version{
			Version: version,
		},
	}

	outputResponse(response)
}

func fatal(message string, err error) {
	_, _ = fmt.Fprintf(os.Stderr, "error %s: %s\n", message, err)
	os.Exit(1)
}

func inputRequest(request *in.Request) {
	if err := json.NewDecoder(os.Stdin).Decode(request); err != nil {
		log.Fatal("[IN] reading request from stdin: ", err)
	}
}

func outputResponse(response in.Response) {
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatal("[IN] writing response to stdout: ", err)
	}
}

