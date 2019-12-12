
FROM docker.pkg.github.com/esportsdrafts/esportsdrafts/esportsdrafts-golang:10-28-2019

COPY go.mod go.sum /workspace/
COPY vendor /workspace/vendor
COPY libs /workspace/libs

ARG VERSION=unknown

LABEL Name "esportsdrafts-base"
LABEL Version ${VERSION}

CMD ["ash"]
