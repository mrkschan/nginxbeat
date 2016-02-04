FROM golang:1.5.3
MAINTAINER KS Chan <mrkschan@gmail.com>

## Install go package dependencies
RUN set -x \
  go get \
	github.com/pierrre/gotestcover \
	github.com/tsg/goautotest \
	golang.org/x/tools/cmd/cover \
	golang.org/x/tools/cmd/vet

ENV GO15VENDOREXPERIMENT=1
