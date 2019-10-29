
FROM esportsdrafts-base:latest AS builder

COPY ./ /workspace/services/auth

RUN cd /workspace/services/auth/cmd/ && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a -ldflags="-w -s" -installsuffix cgo -mod=vendor \
    -o /app .

FROM scratch

ARG VERSION=unknown
ENV VERSION=${VERSION}

LABEL Name "esportsdrafts-auth"
LABEL Version ${VERSION}

COPY --from=builder /etc/ssl/certs/ca-certificates.crt \
    /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app /app

# TODO: Configure app through env variables for more flexibility in k8s
CMD ["/app", "-port", "8000"]
