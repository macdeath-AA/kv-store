# Build stage
FROM golang:1.25.3-alpine AS builder
WORKDIR /build
COPY . .
WORKDIR /build/rest-api
RUN go build -o /build/rest-api-app ./cmd/main.go

# Final image
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /build/rest-api-app .
EXPOSE 8080
CMD ["./rest-api-app"]