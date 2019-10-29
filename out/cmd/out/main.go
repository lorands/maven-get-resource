package main

import (
	"encoding/json"
	"fmt"
	"log"
	// "io/ioutil"

	"github.com/lorands/maven-get-resource"
	"github.com/lorands/maven-get-resource/out"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal(fmt.Sprintf("usage: %v <sources directory>", os.Args[0]))
	}

	//output to stdout...
	response := out.Response{
		Version: resource.Version{
			Version: "N/A",
		},
		Metadata: []resource.MetadataPair {

		},
	}

	outputResponse(response)
}


func outputResponse(response out.Response) {
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatal("[OUT] writing response to stdout", err)
	}
}
