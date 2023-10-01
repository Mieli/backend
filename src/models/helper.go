package models

import "time"

type EventRecord struct {
	CreatedAt *time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt" bson:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt" bson:"deletedAt"`
}
