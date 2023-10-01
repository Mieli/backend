package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	EventRecord *EventRecord       `json:"eventRecord" bson:"eventRecord"`
}

type TaskRepository interface {
	FindAll() ([]Task, error)
	FindById(id string) (*Task, error)
	Add(task Task) error
	Update(id string, task Task) error
	Remove(id string) error
}
