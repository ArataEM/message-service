FROM golang:1.23.3-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . . 

RUN go build -o bin/message-service cmd/main.go



FROM scratch

WORKDIR /

COPY --from=build /app/bin/message-service /message-service

EXPOSE 8080

ENTRYPOINT ["/message-service"]
