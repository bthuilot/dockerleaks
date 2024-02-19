FROM golang:1.20 as build

WORKDIR /build
COPY cmd /build/cmd
COPY internal /build/internal
COPY pkg /build/pkg
COPY go.mod go.sum main.go /build/

RUN GOOS=linux go build -ldflags "-s -w" -o /build/dockerleaks .

FROM alpine:3.18

RUN apk add libc6-compat
RUN rm -rf /sbin/apk

RUN apt update && apt install -y ca-certificates && update-ca-certificates

WORKDIR /app

COPY  --from=build  /build/dockerleaks /app/dockerleaks
COPY  dockerleaks.yml /app/dockerleaks.yml

ARG VERSION="+unknown"

LABEL authors="bryce@thuilot.io"
LABEL repository="github.com/bthuilot/dockerleaks"
LABEL version="${VERSION}"
ENTRYPOINT ["/app/dockerleaks"]