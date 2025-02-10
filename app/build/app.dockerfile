FROM golang:alpine AS build

RUN apk add --no-cache git

RUN mkdir /app

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN go build -o /app/bin ./cmd/app.go

FROM alpine:edge

COPY --from=build /app/bin /app/bin

COPY .env .env

ENTRYPOINT ["/app/bin"]