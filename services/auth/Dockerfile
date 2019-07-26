
FROM efantasy-base:latest AS builder

ENV GOFLAGS=-mod=vendor

RUN apk --no-cache add git ca-certificates

COPY go.mod go.sum /workspace/

RUN go install -v ./...

COPY ./ /workspace

RUN cd cmd/ && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a -ldflags="-w -s" -installsuffix cgo \ 
    -o /app .

FROM scratch

ARG VERSION=unknown
ENV VERSION=${VERSION}

COPY --from=builder /etc/ssl/certs/ca-certificates.crt \
    /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app /app

CMD ["/app", "-port", "8000"]

