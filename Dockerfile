############################################################
# Dockerfile to install Go Monitor
# Based on docker.io/golang Image
############################################################

FROM docker.io/golang:latest

ENV HTTP_PROXY http://127.0.0.1:3128
ENV HTTPS_PROXY http://127.0.0.1:3128

ENV GOPATH /go
ENV PATH /usr/src/go/bin:$PATH
ENV PATH /go/bin:$PATH

#COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

RUN go get github.com/mathcunha/gomonitor
RUN go get gopkg.in/mgo.v2/bson
RUN go get github.com/mathcunha/amon/scheduler

#RUN echo 'package main\nimport (\n\t"github.com/mathcunha/amon"\n)\n\nfunc main(){\n\twg, _ := amon.Monitor("/config/config.json")\n\twg.Wait()\n}' > /go/src/main.go
RUN echo '#!/bin/bash \n unset HTTP_PROXY \n unset HTTPS_PROXY \n go run /go/src/github.com/mathcunha/gomonitor/gomonitor.go -config=/config/config.json' > /entrypoint.sh
RUN chmod 777 /entrypoint.sh


ENTRYPOINT ["/entrypoint.sh"]

# Set the file maintainer (your name - the file's author)
MAINTAINER Matheus Cunha
