
FROM efantasy-base:latest as build_stage

COPY go.mod go.sum /workspace/

RUN go install -v ./...

COPY ./ /workspace

RUN cd cmd/ && \
    CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app .

FROM alpine:3.9

COPY --from=build_stage /app /app

ARG VERSION=unknown

RUN apk --no-cache add tini ca-certificates

EXPOSE 8000

USER nobody

ENTRYPOINT [ "/sbin/tini", "--" ]
CMD ["/app", "-port", "8000"]
