FROM golang:1.26.4-alpine3.24@sha256:3ad57304ad93bbec8548a0437ad9e06a455660655d9af011d58b993f6f615648 AS base
EXPOSE 8080
WORKDIR /src
COPY . .

FROM base AS build
WORKDIR /src
RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o ./bin/app

FROM alpine:3.24.1@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b
WORKDIR /
COPY --from=build /src/bin/app .
ENTRYPOINT ["./app"]
