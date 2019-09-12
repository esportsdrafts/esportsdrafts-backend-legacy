
FROM docker.pkg.github.com/barreyo/efantasy/efantasy-golang:27c1535

ARG VERSION=unknown

LABEL Name "eFantasy-base"
LABEL Version ${VERSION}

COPY go.mod go.sum /workspace/
COPY vendor /workspace/vendor
COPY libs /workspace/libs

CMD ["ash"]
