# syntax=docker/dockerfile:1

## Build
FROM golang:1.19-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY pkg/ ./pkg/

RUN go build -o /docker-keyvalue-service

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /docker-keyvalue-service /docker-keyvalue-service

EXPOSE 5000

USER nonroot:nonroot

ENTRYPOINT ["/docker-keyvalue-service"]