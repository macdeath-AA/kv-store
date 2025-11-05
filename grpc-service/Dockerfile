# Build stage
FROM golang:1.25.3-alpine AS builder
WORKDIR /build
COPY . . 
WORKDIR /build/grpc-service
RUN go build -o /build/grpc-server ./cmd/main.go

# Final image
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /build/grpc-server .
EXPOSE 50051
CMD ["./grpc-server"]