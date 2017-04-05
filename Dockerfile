FROM golang:1.7.5-alpine
MAINTAINER Michael Kraus, michael.kraus@consol.de

RUN apk add --update bash && rm -rf /var/cache/apk/*
RUN mkdir -p /go/src/app
VOLUME /go/src/app
WORKDIR /go/src/app

CMD ["/go/src/app/build.sh"]
