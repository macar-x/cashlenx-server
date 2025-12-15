# syntax=docker/dockerfile:1
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go mod download
ENTRYPOINT ["go", "run", "main.go", "server", "start"]
