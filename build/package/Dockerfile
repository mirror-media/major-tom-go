FROM golang:1.16-alpine3.13 AS build

WORKDIR /major-tom

COPY . .

RUN apk add --update --no-cache make gcc libc-dev && \
    go get ./... && \
    make all

FROM alpine:3.13

WORKDIR /major-tom

COPY --from=build /major-tom/bin/ .

ENTRYPOINT [ "/usr/bin/env", "./bot" ]
