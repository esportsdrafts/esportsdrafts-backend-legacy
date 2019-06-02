
FROM golang:1.12-alpine3.9 AS BUILD

ONBUILD RUN go build -v -o /run_service ./workspace/...

FROM alpine:3.9

ENV WORKSPACE /workspace
RUN mkdir /workspace
WORKDIR /workspace

ONBUILD COPY --from=BUILD /run_service /workspace

RUN apk --no-cache add tini ca-certificates

ENTRYPOINT [ "/sbin/tini", "--" ]
