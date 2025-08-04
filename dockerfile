# Dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /gobank

FROM alpine:latest
WORKDIR /
COPY --from=builder /gobank /gobank
EXPOSE 3000
CMD ["/gobank"]
