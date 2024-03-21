ARG GO_VERSION="1.20"
ARG RUNNER_IMAGE="golang:${GO_VERSION}-alpine3.17"

FROM golang:${GO_VERSION}-bullseye AS build-env

WORKDIR /go/src/github.com/sagaxyz/sagacli

RUN apt-get update -y

COPY . .

RUN make build

FROM ${RUNNER_IMAGE}

COPY --from=build-env /go/src/github.com/sagaxyz/sagacli/build/sscd /usr/bin/sscd

RUN apk add gcompat bash

EXPOSE 26656
EXPOSE 26660
EXPOSE 26657
EXPOSE 1317
EXPOSE 9090

CMD ["sscd", "start"]