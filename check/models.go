package check

import resource "github.com/lorands/maven-get-resource"

type Request struct {
	Source  resource.Source  `json:"source"`
	Version resource.Version `json:"version"`
}

type Response []resource.Version
