FROM golang:1.16-alpine3.14 AS build

WORKDIR /major-tom

COPY . .

RUN apk add --update --no-cache make gcc libc-dev && \
    go mod download && \
    make all

FROM alpine:3.14

WORKDIR /major-tom

COPY --from=build /major-tom/bin/ .

ENTRYPOINT [ "/usr/bin/env", "./bot" ]
