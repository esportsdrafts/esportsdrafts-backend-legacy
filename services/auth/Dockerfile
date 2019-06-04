
FROM efantasy-base:latest as build_stage

COPY ./ /workspace

RUN cd auth && \
    go get -u ./...
RUN cd api && \
    go install -v ./...

RUN cd api/ && \
    CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app .

FROM alpine:3.9

COPY --from=build_stage /app /app

RUN apk --no-cache add tini ca-certificates

EXPOSE 8000

USER nobody

ENTRYPOINT [ "/sbin/tini", "--" ]
CMD ["/app"]
