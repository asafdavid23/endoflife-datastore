package mongo

import (
	"context"
	"log"
	"time"

	"github.com/asafdavid23/endoflife-datastore/internal/api"
	"github.com/asafdavid23/endoflife-datastore/internal/models"
	v1 "github.com/asafdavid23/endoflife-operator/api/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect establishes a connection to the MongoDB server.
func Connect(ctx context.Context, uri string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the MongoDB server to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}

// UpdateMongoDB updates the MongoDB collection with ProductCheck data.
func UpdateMongoDB(ctx context.Context, mongoCollection *mongo.Collection, productCheck v1.ProductCheck) error {
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

func CheckAndNotify(collection *mongo.Collection, timeframe time.Duration) {
	now := time.Now()
	thresholdDate := now.Add(timeframe)

	// Query for products with EOL dates within the specified timeframe
	filter := bson.M{
		"endOfLifeDate": bson.M{
			"$gte": now,
			"$lte": thresholdDate,
		},
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		return
	}
	defer cursor.Close(context.Background())

	var products []models.Product
	if err = cursor.All(context.Background(), &products); err != nil {
		log.Printf("Error decoding products: %v", err)
		return
	}

	if len(products) > 0 {
		// Assuming you have a function sendNotifications in the current package
		api.SendNotifications(products)
		log.Printf("Sent notifications for %v products nearing end-of-life.", len(products))
	} else {
		log.Println("No products nearing end-of-life within the specified timeframe.")
	}
}
