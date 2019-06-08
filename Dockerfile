
FROM golang:1.12-alpine3.9

ENV GO111MODULE=on

ENV WORKSPACE /workspace
RUN mkdir /workspace
WORKDIR /workspace
