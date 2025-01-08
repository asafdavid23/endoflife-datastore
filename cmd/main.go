package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/asafdavid23/endoflife-datastore/internal/config"
	"github.com/asafdavid23/endoflife-datastore/internal/k8s"
	"github.com/asafdavid23/endoflife-datastore/internal/mongo"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	// Load configuration
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Loaded configuration: %+v", cfg)

	// Set up in-cluster Kubernetes client
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to load in-cluster Kubernetes config: %v", err)
	}

	k8sClient, err := client.New(restConfig, client.Options{})
	if err != nil {
		log.Fatalf("Failed to initialize Kubernetes client: %v", err)
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, os.Getenv("MONGODB_URI"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	mongoCollection := mongoClient.Database(os.Getenv("MONGODB_DATABASE")).Collection(os.Getenv("MONGODB_COLLECTION"))
	log.Printf("Connected to MongoDB collection: %s", os.Getenv("MONGODB_COLLECTION"))

	// Watch and process ProductCheck objects
	if err := k8s.WatchAndProcessProductChecks(ctx, k8sClient, mongoCollection, cfg.Kubernetes.Namespace); err != nil {
		log.Fatalf("Error watching ProductCheck objects: %v", err)
	}
}
