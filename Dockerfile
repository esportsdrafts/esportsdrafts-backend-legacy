
FROM docker.pkg.github.com/barreyo/esportsdrafts/esportsdrafts-golang:10-28-2019

ARG VERSION=unknown

LABEL Name "esportsdrafts-base"
LABEL Version ${VERSION}

COPY go.mod go.sum /workspace/
COPY vendor /workspace/vendor
COPY libs /workspace/libs

CMD ["ash"]
