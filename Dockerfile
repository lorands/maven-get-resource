FROM golang:alpine as builder
#RUN apk add --no-cache curl jq
RUN mkdir -p /assets
WORKDIR /assets
COPY . /go/src/github.com/lorands/maven-get-resource
COPY ./assets /assets
ENV CGO_ENABLED 0
RUN go build -o /assets/in github.com/lorands/maven-get-resource/in/cmd/in
RUN go build -o /assets/out github.com/lorands/maven-get-resource/out/cmd/out
RUN go build -o /assets/check github.com/lorands/maven-get-resource/check/cmd/check
#RUN mkdir -p /tests
#RUN go test -o /tests/ -c github.com/lorands/maven-get-resource
WORKDIR /go/src/github.com/lorands/maven-get-resource
#RUN set -e; for pkg in $(go list ./... | grep -v "acceptance"); do \
#		go test -o "/tests/$(basename $pkg).test" -c $pkg; \
#	done
#RUN go test .

FROM openjdk:8u151-jdk-alpine AS resource
RUN apk add --no-cache bash tzdata ca-certificates
# https://github.com/concourse/concourse/issues/2042
RUN unlink  $JAVA_HOME/jre/lib/security/cacerts && \
    cp /etc/ssl/certs/java/cacerts $JAVA_HOME/jre/lib/security/cacerts
COPY --from=builder assets/ /opt/resource/
RUN chmod +x /opt/resource/*

#FROM resource AS tests
#COPY --from=builder /tests /go-tests
#COPY out/assets /go-tests/assets
#WORKDIR /go-tests
#RUN set -e; for test in /go-tests/*.test; do \
#    $test; \
#done

FROM resource
