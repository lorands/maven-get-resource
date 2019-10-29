package maven

import "encoding/xml"

type MavenMetadata struct {
	XMLName      xml.Name   `xml:"metadata"`
	ModelVersion string     `xml:"modelVersion,attr,omitempty"`
	GroupID      string     `xml:"groupId"`
	ArtifactID   string     `xml:"artifactId"`
	Versioning   Versioning `xml:"versioning"`
}

type Versioning struct {
	XMLName     xml.Name `xml:"versioning"`
	Latest      string   `xml:"latest,omitempty"`
	Release     string   `xml:"release"`
	Versions    Versions `xml:"versions"`
	LastUpdated string   `xml:"lastUpdated"`
}

type Versions struct {
	XMLName xml.Name  `xml:"versions"`
	Version []string `xml:"version"`
}
