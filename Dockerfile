# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o urlshortener

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/urlshortener .
COPY templates/ templates/
COPY static/ static/

EXPOSE 8080
CMD ["./urlshortener"]
