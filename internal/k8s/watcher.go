package k8s

import (
	"context"
	"os"
	"time"

	"github.com/asafdavid23/endoflife-datastore/internal/logging"
	"github.com/asafdavid23/endoflife-datastore/internal/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// WatchAndProcessProductChecks watches ProductCheck objects and updates MongoDB.
func WatchAndProcessProductChecks(ctx context.Context, k8sClient client.Client, mongoCollection *mongo.Collection, namespace string) error {
	logLevel := os.Getenv("LOG_LEVEL")
	logger := logging.NewLogger(logLevel)

	logger.Printf("Starting to watch ProductCheck objects in namespace: %s", namespace)

	for {
		// Fetch ProductCheck objects
		productChecks, err := FetchProductChecks(ctx, k8sClient, namespace)
		if err != nil {
			logger.Printf("Failed to fetch ProductCheck objects: %v", err)
			return err
		}

		// Process each ProductCheck
		for _, productCheck := range productChecks {
			if err := UpdateMongoDB(ctx, mongoCollection, productCheck); err != nil {
				logger.Printf("Failed to update MongoDB for ProductCheck %s: %v", productCheck.Name, err)
			} else {
				logger.Printf("Successfully updated MongoDB for ProductCheck %s", productCheck.Name)
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
func FetchProductChecks(ctx context.Context, k8sClient client.Client, namespace string) ([]models.ProductCheck, error) {
	productCheckList := &models.ProductCheckList{}
	err := k8sClient.List(ctx, productCheckList, &client.ListOptions{Namespace: namespace})
	if err != nil {
		return nil, err
	}

	return productCheckList.Items, nil
}

// UpdateMongoDB updates the MongoDB collection with ProductCheck data.
func UpdateMongoDB(ctx context.Context, mongoCollection *mongo.Collection, productCheck models.ProductCheck) error {
	filter := map[string]interface{}{"name": productCheck.Name}
	update := map[string]interface{}{
		"$set": map[string]interface{}{
			"name":       productCheck.Name,
			"status":     productCheck.Status,
			"lastUpdate": time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := mongoCollection.UpdateOne(ctx, filter, update, opts)
	return err
}
