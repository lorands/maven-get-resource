# maven-get-resource

Maven (m2) Concourse get resource.

## Source Configuration

* `url`: *Required*. The target maven repository location.
* `artifact`: *Required*. The artifact coordinates in the form of groupId:artifactId:type\[:classifier\]
* `username`: *Optional*. The username used to authenticate.
* `password`: *Optional*. The password used to authenticate.
* `verbose`: *Optional*. True to write intensive log.

## Check: Check if there is new version in source repository

Check if there is new version of the artifact

## Get (in): Download the artifact and deploy it to target repository

Downloads the artifact, pom file, sha1, etc..
 
It will create the following files in the target directory:

- pom file
- archive file (usually jar)
- version file 

## Put (out): Nothing.

Does nothing.

## Pipeline example

```yaml
resource_types:
  - name: maven-get-resource
    type: docker-image
    source:
      repository: lorands/maven-get-resource
resources:
  - name: dev-artifact
    type: maven-get-resource
    source:
      url: https://mynexus.example.com/repository/develop
      artifact: my.group:my-artifact:jar
      username: myUser
      password: myPass
jobs:
  - name: myJob
    plan:
    - get: dev-artifact
```







