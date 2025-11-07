package main

import (
	"context"
	"log"
	"net/http"
	"time"

	pb "kv-store/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/gin-gonic/gin"
)

// gRPC connection info
const grpcAddr = "grpc-service:50051"

func main() {
	router := gin.Default()

	// connect to gRPC backend
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC backend: %v", err)
	}
	defer conn.Close()
	client := pb.NewKVStoreClient(conn)

	// set key-value
	router.POST("/kv", func(c *gin.Context) {
		var req struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		resp, err := client.Set(ctx, &pb.SetRequest{Key: req.Key, Value: req.Value})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": resp.Success, "message": resp.Message})
	})

	// get value by key
	router.GET("/kv/:key", func(c *gin.Context) {
		key := c.Param("key")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		resp, err := client.Get(ctx, &pb.GetRequest{Key: key})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"key": key, "value": resp.Value, "found": resp.Found})
	})

	// delete key
	router.DELETE("/kv/:key", func(c *gin.Context) {
		key := c.Param("key")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		resp, err := client.Delete(ctx, &pb.DeleteRequest{Key: key})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": resp.Success})
	})

	// start server
	log.Println("REST API listening on :8080")
	router.Run(":8080")
}