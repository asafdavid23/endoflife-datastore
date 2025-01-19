package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/asafdavid23/endoflife-datastore/internal/config"
	"github.com/asafdavid23/endoflife-datastore/internal/k8s"
	"github.com/asafdavid23/endoflife-datastore/internal/mongo"
	v1 "github.com/asafdavid23/endoflife-operator/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
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

	log.Printf("Loaded env vars : %+v", os.Environ())
	log.Printf("Loaded configuration: %+v", cfg)

	// Create a scheme and register the CRD types
	scheme := runtime.NewScheme()
	utilruntime.Must(v1.AddToScheme(scheme))             // Register ProductCheck and ProductCheckList
	utilruntime.Must(clientgoscheme.AddToScheme(scheme)) // Register built-in types

	// Set up in-cluster Kubernetes client
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to load in-cluster Kubernetes config: %v", err)
	}

	k8sClient, err := client.New(restConfig, client.Options{Scheme: scheme})
	if err != nil {
		panic(err)
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, os.Getenv("MONGODB_URI"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	mongoCollection := mongoClient.Database(os.Getenv("MONGODB_DATABASE")).Collection(os.Getenv("MONGODB_COLLECTION"))
	log.Printf("Connected to MongoDB collection: %s", os.Getenv("MONGODB_COLLECTION"))

	// Watch and process ProductCheck objects
	if err := k8s.WatchAndProcessProductChecks(ctx, k8sClient, mongoCollection, os.Getenv("WATCH_NAMESPACE")); err != nil {
		log.Fatalf("Error watching ProductCheck objects: %v", err)
	}

	// Check and notify for EOL products
	timeframe := -30 * 24 * time.Hour

	// Ticker to run the checkAndNotify function every 24 hours
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mongo.CheckAndNotify(mongoCollection, timeframe)
		}
	}
}
