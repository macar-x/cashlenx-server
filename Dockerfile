FROM golang:1.21-alpine
WORKDIR /app

# Install curl for health checks
RUN apk add --no-cache curl

COPY . .
RUN go mod download
ENTRYPOINT ["go", "run", "main.go", "server", "start"]
