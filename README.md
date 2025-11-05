# Decomposed Key-Value Store
This project implements a simple distributed Key-Value store in Go, featuring two services, Rest API and gRPC, which communicate over gRPC.

The Key-Value store service can:
- Store a value at a given key
- Retrieve the value for a given key
- Delete a given key

The application is set up using `docker-compose`.

## Running
```bash
  docker-compose up --build
```
Both containers (`rest-api` and `grpc-service`) should print logs showing they are ready.

## API Endpoints
- `POST /kv` - Set a key: `{ "key": "hello", "value": "world" }`
- `GET /kv/:key` - Get the value at key
- `DELETE /kv/:key` - Delete key

## Testing
1. Store a new key-value pair (POST)
```bash
  curl -X POST http://localhost:8080/kv -H "Content-Type: application/json" -d '{"key":"city","value":"newyork"}'
```
2. Retrieve a value for a given key (GET)
```bash
  curl http://localhost:8080/kv/name
```
3.  Delete a key (DELETE):
```bash
curl -X DELETE http://localhost:8080/kv/name
```
Invalid or missing keys are handled appropriately. Each service is a standalone Go module with its own `go.mod`. 
