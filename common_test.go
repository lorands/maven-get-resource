package resource

import (
	//"io/ioutil"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestSource_GetMavenMetadata(t *testing.T) {
	source := Source{
		URL:      "https://repo1.maven.org/maven2",
		Artifact: "org.apache.commons:commons-lang3",
		Username: "",
		Password: "",
		Verbose:  true,
	}

	metadata, e := source.GetMavenMetadata()
	if e != nil {
		t.Errorf("Fail to get metadata: %v", e)
	}

	if len(metadata.Versioning.Versions.Version) <= 0 {
		t.Errorf("No version information!")
	}

	if len(metadata.Versioning.Release) <= 0 {
		t.Errorf("Missing release info.")
	}
}

func TestInResource_Download(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "gogo")
	pkg := reflect.TypeOf(InResource{}).PkgPath()

	assets := filepath.Join(findExsitingGoPkgPath(pkg), "assets")

	in := InResource{
		Source: Source{
			URL:      "https://repo1.maven.org/maven2",
			Artifact: "org.apache.commons:commons-lang3",
			Username: "",
			Password: "",
			Verbose:  true,
		},
		Version: Version{
			Version: "3.8.1",
		},
		DestinationDir: tmp,
		ResourceDir:    assets,
	}

	err := in.Download()
	if err != nil {
		t.Error(err)
	}

	//check if files are there...
	files, _ := ioutil.ReadDir(tmp)

	if len(files) < 3 {
		t.Errorf("Excpected to have at least 3 files, insted we have %d files", len(files))
	}

	cntr := 0
	for _, file := range files {
		if file.Name() == "commons-lang3-3.8.1.jar" {
			cntr++
		} else if file.Name() == "commons-lang3-3.8.1.pom" {
			cntr++
		} else if file.Name() == "version" {
			cntr++
		} else {
			//t.Errorf("Unknow file found: %s", file.Name())
		}
	}
	if cntr < 3 {
		t.Errorf("Not all required files found, expecting 3 well defined, found: %d", cntr)
	}

}

func findExsitingGoPkgPath(pkg string) string {
	goPaths := strings.Split(os.Getenv("GOPATH"), string(os.PathListSeparator))
	for _, path := range goPaths {
		pkgPath := filepath.Join(path, "src", pkg)
		info, _ := os.Stat(pkgPath)
		if info.IsDir() {
			return pkgPath
		}
	}

	return ""
}