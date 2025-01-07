package api

import (
	"context"
	"log"
	"net/http"

	"github.com/asafdavid23/endoflife-datastore/internal/k8s"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Server defines the HTTP server.
type Server struct {
	router          *gin.Engine
	k8sClient       client.Client
	mongoCollection *mongo.Collection
	namespace       string
}

// NewServer creates a new API server instance.
func NewServer(k8sClient client.Client, mongoCollection *mongo.Collection, namespace string) *Server {
	s := &Server{
		router:          gin.Default(),
		k8sClient:       k8sClient,
		mongoCollection: mongoCollection,
		namespace:       namespace,
	}
	s.setupRoutes()
	return s
}

// setupRoutes sets up the HTTP API endpoints.
func (s *Server) setupRoutes() {
	s.router.GET("/test", s.testHandler)
}

// Start runs the API server on the specified port.
func (s *Server) Start(port string) error {
	log.Printf("Starting API server on port %s", port)
	return s.router.Run(":" + port)
}

// testHandler handles the /test endpoint to test Kubernetes and MongoDB interactions.
func (s *Server) testHandler(c *gin.Context) {
	ctx := context.Background()

	// Fetch ProductCheck objects
	productChecks, err := k8s.FetchProductChecks(ctx, s.k8sClient, s.namespace)
	if err != nil {
		log.Printf("Failed to fetch ProductCheck objects: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ProductChecks"})
		return
	}

	// Simulate updating MongoDB
	for _, productCheck := range productChecks {
		err := k8s.UpdateMongoDB(ctx, s.mongoCollection, productCheck)
		if err != nil {
			log.Printf("Failed to update MongoDB for ProductCheck %s: %v", productCheck.Name, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update MongoDB"})
			return
		}
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{
		"message":        "ProductChecks processed successfully",
		"processedCount": len(productChecks),
	})
}
