FROM golang:1.19.0-bullseye AS build-env
ARG GITHUB_USER
ARG GITHUB_TOKEN

WORKDIR /root

RUN apt-get update -y
RUN apt-get install git jq -y

COPY . .

RUN git config --global --add url."https://${GITHUB_USER}:${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"
RUN go env -w GOPRIVATE='github.com/sagaxyz/*'
RUN go build -o /root/ ./...
RUN cp /root/sscd /usr/bin/sscd

FROM golang:1.19.0-bullseye

RUN apt-get update -y
RUN apt-get install ca-certificates jq -y
RUN apt-get install -y pip
RUN pip install --upgrade awscli

WORKDIR /root

COPY --from=build-env /root/sscd /usr/bin/sscd
COPY --from=build-env /root/start.sh /root/start.sh
COPY --from=build-env /root/defaults.genesis.json /root/defaults.genesis.json

RUN chmod -R 755 /root/start.sh
RUN /usr/local/bin/aws --version

EXPOSE 26656 26657 1317 9090

CMD ["bash","/root/start.sh"]
