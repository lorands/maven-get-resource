package resource

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/lorands/maven-get-resource/maven"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
)


func (in InResource) Download() error {
	//archivePath := filepath.Join(srcDir, archiveFileName)
	//pomPath := filepath.Join(srcDir, pomFileName)

	gav := fmt.Sprintf("%s:%s",in.Source.Artifact, in.Version.Version)

	args := [] string {
		"-s",
		"settings.xml",
		"org.apache.maven.plugins:maven-dependency-plugin:3.1.1:get",
		fmt.Sprintf("-Drepository.url=%s", in.Source.URL),
		fmt.Sprintf("-Drepository.username=%s", in.Source.Username),
		fmt.Sprintf("-Drepository.password=%s", in.Source.Password),
		fmt.Sprintf("-Dartifact=%s", gav),
		"-Dtransitive=false",
	}

	if err := runCmd(in.ResourceDir, "./mvnw", args); err != nil {
		return err
	}

	//copy to dest
	homeDir, err := os.UserHomeDir()
	if err != nil  {
		return err
	}
	artifactDef, err := ArtifactStrToArtifactDef(in.Source.Artifact)
	if err != nil  {
		return nil
	}

	sourceDir := filepath.Join(homeDir, ".m2", "repository", artifactDef.toRelativePath(), in.Version.Version)

	if err := CopyDirectory(sourceDir, in.DestinationDir); err != nil {
		return err
	}
	//write version file
	versionFile := filepath.Join(in.DestinationDir, "version")
	if fileExists(versionFile) {
		if in.Source.Verbose {
			_, _ = fmt.Fprintf(os.Stderr, "Version file already exists with version name. Skipping...")
		}
	} else {
		if err = ioutil.WriteFile(versionFile, []byte(in.Version.Version), 0644); err != nil {
			return err
		}
	}
	return nil
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func runCmd(workDir string, cmdStr string, strParams []string) error {
	//tracelog(trace,"About to Execute %s with params: %v", cmdStr, strParams)
	cmd := exec.Command(cmdStr, strParams...)
	cmd.Dir = workDir
	var sout bytes.Buffer
	cmd.Stdout = &sout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() != 0 {
					return fmt.Errorf("Non zero exit code form cli: %d\n%s", status.ExitStatus(), sout.String())
				}
			}
		} else {
			//log.Fatalf("cmd.Wait: %v", err)
			return err
		}
	}
	return nil
}


func (source Source) GetMavenMetadata() (maven.MavenMetadata, error) {
	var result maven.MavenMetadata

	r := regexp.MustCompile("\\/$")
	srcURL := r.ReplaceAllString(source.URL, "")

	adef, err := ArtifactStrToArtifactDef(source.Artifact)
	if err != nil {
		return result, fmt.Errorf("fail to process artfiact from resource source. %v", err)
	}
	groupPath := strings.Replace(adef.GroupID, ".", "/", -1)
	metaURL := strings.Join([]string{srcURL, groupPath, adef.ArtifactID, "maven-metadata.xml"}, "/")
	var client http.Client
	req, err := http.NewRequest("GET", metaURL, nil)
	if err != nil {
		return result, fmt.Errorf("fail to create request object to maven-metadata.xml: %v", err)
	}
	if source.Username != "" {
		//main2.tracelog("Setting basic authorization as requested for user: %v\n", username)
		req.SetBasicAuth(source.Username, source.Password)
	}
	resp, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("Error response from http. %v\n%v", resp, err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		if source.Verbose {
			_, _ = fmt.Fprintf(os.Stderr, "Status code for download: %d\n", resp.StatusCode)
		}
	} else {
		return result, fmt.Errorf("fail to download artifact. Status code %s", resp.Status)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("fail to read maven-metadata.xml: %v", err)
	}
	err = xml.Unmarshal(data, &result)
	if err != nil {
		return result, fmt.Errorf("fail to process xml: %v", err)
	}
	return result, nil
}

func (artifactDef ArtifactDef) toRelativePath() string {
	parts := strings.Split(artifactDef.GroupID, ".")
	parts = append(parts, artifactDef.ArtifactID)
	return filepath.Join(parts...)
}

func ArtifactStrToArtifactDef(artifact string) (ArtifactDef, error) {
	var def ArtifactDef

	splits := strings.Split(artifact, ":")

	if len(splits) < 2 {
		err := fmt.Errorf("you must specify at least GroupID:ArtifactID")
		return def, err
	}

	def.GroupID = splits[0]
	def.ArtifactID = splits[1]
	//def.AType = splits[2]

	//if len(splits) > 3 {
	//	def.Classifier = splits[3]
	//}

	return def, nil
}