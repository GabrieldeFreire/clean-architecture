# Dockerfile
FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
COPY cmd/ordersystem/.env .

RUN go build -tags wireinject -o main cmd/ordersystem/main.go cmd/ordersystem/wire_gen.go
CMD [ "./main" ]

