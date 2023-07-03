FROM golang:1.20 as build

WORKDIR /build
COPY cmd /build/cmd
COPY internal /build/internal
COPY pkg /build/pkg
COPY go.mod go.sum main.go /build/

RUN GOOS=linux go build -ldflags "-s -w" -o /build/dockerleaks .

FROM debian:bookworm-slim

WORKDIR /app
COPY --from=build /build/dockerleaks /app/dockerleaks

LABEL authors="bryce@thuilot.io"
LABEL repository="github.com/bthuilot/dockerleaks"
ENTRYPOINT ["/app/dockerleaks"]