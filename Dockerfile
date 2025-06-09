ARG GO_VERSION="1.24.4"
FROM golang:${GO_VERSION}-bookworm AS build-env
ARG GITHUB_USER
ARG GITHUB_TOKEN

WORKDIR /root

RUN apt-get update -y

COPY . .

RUN git config --global --add url."https://${GITHUB_USER}:${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"
RUN make build

FROM golang:${GO_VERSION}-alpine3.22

COPY --from=build-env /root/build/sscd /usr/bin/
COPY --from=build-env /root/start.sh /root/

RUN apk add gcompat bash curl

EXPOSE 26656
EXPOSE 26660
EXPOSE 26657
EXPOSE 1317
EXPOSE 9090

CMD ["bash","/root/start.sh"]
