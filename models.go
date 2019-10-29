package resource

type Source struct {
	URL         string `json:"source_url"`
	Artifact    string `json:"artifact"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Verbose     bool   `json:"verbose,omitempty"`
}

type Version struct {
	Version string `json:"version,omitempty"`
}

type MetadataPair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ---

//result, err := resource.GetMavenMetadata(request.Source.Artifact, srcURL, request.Source.Username, request.Source.Password)

//type MetadataRequest struct {
//	Source Source
//}

type InResource struct {
	Source Source
	Version Version
	DestinationDir string
	ResourceDir string
}

// ArtifactDef defines an artifact as a struct.
type ArtifactDef struct {
	GroupID    string
	ArtifactID string
	//AType      string
	//Classifier string
}
