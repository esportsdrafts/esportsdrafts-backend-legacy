
FROM efantasy-base:latest

ENV GOFLAGS=-mod=vendor

RUN apk --no-cache add git tini ca-certificates

COPY go.mod go.sum /workspace/

RUN go install -v ./...

COPY ./ /workspace

RUN cd cmd/ && \
    CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app .

ARG VERSION=unknown
ENV VERSION=${VERSION}

EXPOSE 8000

USER nobody

ENTRYPOINT [ "/sbin/tini", "--" ]
CMD ["/app", "-port", "8000"]
