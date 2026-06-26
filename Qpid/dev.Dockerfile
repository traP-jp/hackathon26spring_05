FROM golang:1.26.4-alpine3.24@sha256:3ad57304ad93bbec8548a0437ad9e06a455660655d9af011d58b993f6f615648 AS base
EXPOSE 8080
RUN go install github.com/air-verse/air@v1.65.3
WORKDIR /src
COPY . .

FROM base AS build
RUN go mod download
ENTRYPOINT ["air", "-c", ".air.toml"]
