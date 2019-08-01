
FROM alpine:3.10.1

ARG VERSION="unknown"

LABEL Name="efantasy-beanstalkd"
LABEL Version=${VERSION}

RUN apk add --no-cache beanstalkd tini

EXPOSE 11300

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/usr/bin/beanstalkd"]
