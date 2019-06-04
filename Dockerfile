
FROM golang:1.12-alpine3.9

ENV GO111MODULE=on

RUN apk --no-cache add git

ENV WORKSPACE /workspace
RUN mkdir /workspace
WORKDIR /workspace
