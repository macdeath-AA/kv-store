package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "kv-store/proto"
	"google.golang.org/grpc"
)

// Server struct implements the KVStore service
type server struct {
	pb.UnimplementedKVStoreServer
	store map[string]string
	mu    sync.RWMutex // Thread-safe map access
}

// Set stores a key-value pair
func (s *server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[req.Key] = req.Value
	log.Printf("SET: %s = %s", req.Key, req.Value)

	return &pb.SetResponse{
		Success: true,
		Message: "Key set successfully",
	}, nil
}

// Get retrieves a value for a given key
func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, found := s.store[req.Key]
	log.Printf("GET: %s (found: %v)", req.Key, found)

	return &pb.GetResponse{
		Value: value,
		Found: found,
	}, nil
}

// Delete removes a key from the store
func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.store[req.Key]
	if exists {
		delete(s.store, req.Key)
		log.Printf("DELETE: %s (success)", req.Key)
		return &pb.DeleteResponse{Success: true}, nil
	}

	log.Printf("DELETE: %s (key not found)", req.Key)
	return &pb.DeleteResponse{Success: false}, nil
}

func main() {
	// Create TCP listener on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Create our KV store server with empty map
	kvServer := &server{
		store: make(map[string]string),
	}

	// Register our service
	pb.RegisterKVStoreServer(grpcServer, kvServer)

	fmt.Println("gRPC KV Store server listening on :50051")
	
	// Start serving
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}