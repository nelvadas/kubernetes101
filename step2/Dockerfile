FROM golang:latest
#GOPATH is set to /go in this base image.

ADD .  /go/src/golang/

RUN go install golang/hello

ENTRYPOINT /go/bin/hello

EXPOSE 8080
