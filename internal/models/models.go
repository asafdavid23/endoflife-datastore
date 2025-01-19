package models

import "time"

// Product represents the structure of your product data
type Product struct {
	Name          string    `bson:"name" json:"name"`
	Version       string    `bson:"version" json:"version"`
	EndOfLifeDate time.Time `bson:"endOfLifeDate" json:"endOfLifeDate"`
}
