package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email       string             `json:"email" bson:"email"`
	Password    string             `json:"password" bson:"password"`
	EventRecord *EventRecord       `json:"eventRecord" bson:"eventRecord"`
}

type UserPresenter struct {
	ID    primitive.ObjectID
	Email string
}

type UserRepository interface {
	FindAll() ([]User, error)
	FindById(id string) (*User, error)
	Add(user User) error
	Update(id string, user User) error
	Remove(id string) error
	FindByEmail(email string) (User, error)
	VerifyPassword(user *User, password string) (bool, error)
}
