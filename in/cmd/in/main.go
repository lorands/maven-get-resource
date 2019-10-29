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

var trace bool

func main() {
	var request in.Request

	destinationDir := os.Args[1]
	resourceDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		tracelog("Fail to get resource directory! Bailing out.")
		os.Exit(2)
	}

	tracelog("===> IN!")

	inputRequest(&request)

	trace = request.Source.Verbose
	version := request.Version.Version

	inResource := resource.InResource{
		Source:         resource.Source{},
		Version:        resource.Version{},
		DestinationDir: destinationDir,
		ResourceDir:    resourceDir,
	}

	//version = resource.readIfFile(version)
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
	//reader := bufio.NewReader(os.Stdin)
	//text, _ := reader.ReadString('\n')
	//tracelog("IN: stdin: %s\n", text)
	if err := json.NewDecoder(os.Stdin).Decode(request); err != nil {
	//if err := json.Unmarshal([]byte(text), request); err != nil {
		log.Fatal("[IN] reading request from stdin: ", err)
	}
}

func outputResponse(response in.Response) {
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatal("[IN] writing response to stdout: ", err)
	}
}

func tracelog(message string, args ...interface{}) {
	if trace {
		fmt.Fprintf(os.Stderr, message, args...)
	}
}
