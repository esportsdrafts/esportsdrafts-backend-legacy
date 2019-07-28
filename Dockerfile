
FROM golang:1.13-rc-alpine3.10

ARG VERSION=unknown

LABEL Name "eFantasy-base"
LABEL Version ${VERSION}

ENV WORKSPACE /workspace
RUN mkdir /workspace
WORKDIR /workspace

RUN apk --no-cache add git ca-certificates

COPY go.mod go.sum /workspace/
COPY vendor /workspace/vendor
COPY libs /workspace/libs

CMD ["ash"]
