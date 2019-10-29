package main

import (
	"encoding/json"
	"fmt"
	resource "github.com/lorands/maven-get-resource"
	"github.com/lorands/maven-get-resource/check"
	"log"
	"os"
)

func main() {
	var request check.Request
	inputRequest(&request)

	result, err := request.Source.GetMavenMetadata()
	if err != nil {
		fatal("Fail to get maven metadata and process it.", err)
	}

	var response check.Response

	if result.Versioning.Release != request.Version.Version {
		// for _, version := range result.Versioning.Versions {
		// 	if version.Version < request.Version.Version {
		// 		greaterVersions = append(greaterVersions, version.Version)
		// 	}
		// }
		//releaseVersion := check.Response.Version{Version: "1.1"}
		item := resource.Version{Version: result.Versioning.Release}
		response = append(response, item)
	}

	outputResponse(response)
}

func fatal(message string, err error) {
	_, _ = fmt.Fprintf(os.Stderr, "error %s: %s\n", message, err)
	os.Exit(1)
}

func inputRequest(request *check.Request) {
	if err := json.NewDecoder(os.Stdin).Decode(request); err != nil {
		log.Fatal("[CHK] reading request from stdin", err)
	}
}

func outputResponse(response check.Response) {
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatal("[CHK] writing response to stdout", err)
	}
}