package k8s

import (
	"context"
	"os"
	"time"

	"github.com/asafdavid23/endoflife-datastore/internal/logging"
	v1 "github.com/asafdavid23/endoflife-operator/api/v1"

	internalMongo "github.com/asafdavid23/endoflife-datastore/internal/mongo"
	"go.mongodb.org/mongo-driver/mongo"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// WatchAndProcessProductChecks watches ProductCheck objects and updates MongoDB.
func WatchAndProcessProductChecks(ctx context.Context, k8sClient client.Client, mongoCollection *mongo.Collection, namespace string) error {
	logLevel := os.Getenv("LOG_LEVEL")
	logger := logging.NewLogger(logLevel)
	startFetch := time.Now()

	logger.Printf("FetchProductChecks took %v", time.Since(startFetch))

	for {
		// Fetch ProductCheck objects
		productChecks, err := FetchProductChecks(ctx, k8sClient, namespace)
		if err != nil {
			logger.Printf("Failed to fetch ProductCheck objects: %v", err)
			return err
		}

		// Process each ProductCheck
		for _, productCheck := range productChecks {
			startUpdate := time.Now()
			if err := internalMongo.UpdateMongoDB(ctx, mongoCollection, productCheck); err != nil {
				logger.Fatalf("Failed to update MongoDB for ProductCheck %s: %v", productCheck.Name, err)
			} else {
				logger.Printf("UpdateMongoDB for ProductCheck %s took %v", productCheck.Name, time.Since(startUpdate))
			}
		}

		// Re-run the loop after a delay
		logger.Println("Sleeping before the next watch iteration...")
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(30 * time.Second): // 30-second interval between iterations
		}
	}
}

// FetchProductChecks retrieves ProductCheck objects from Kubernetes.
func FetchProductChecks(ctx context.Context, k8sClient client.Client, namespace string) ([]v1.ProductCheck, error) {
	productCheckList := &v1.ProductCheckList{}
	err := k8sClient.List(ctx, productCheckList, &client.ListOptions{Namespace: namespace})
	if err != nil {
		return nil, err
	}

	return productCheckList.Items, nil
}
