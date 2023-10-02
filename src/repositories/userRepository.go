package repositories

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	pkgmodels "github.com/mieli/backend/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	DB         *mongo.Client
	Collection *mongo.Collection
}

func NewUserRepository(db *mongo.Client) *UserRepository {
	return &UserRepository{
		DB:         db,
		Collection: db.Database(os.Getenv("MONGO_DATABASE")).Collection("users"),
	}
}

func (r *UserRepository) FindAll() ([]pkgmodels.User, error) {
	var users []pkgmodels.User
	ctx := context.TODO()
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Erro ao buscar usuários:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user pkgmodels.User
		if err := cursor.Decode(&user); err != nil {
			log.Println("Erro ao decodificar o usuário:", err)
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
func (r *UserRepository) FindById(id string) (*pkgmodels.User, error) {
	var user pkgmodels.User

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Erro ao converter ID em ObjectID:", err)
		return nil, err
	}
	filter := bson.M{"_id": objectID}
	ctx := context.TODO()
	result := r.Collection.FindOne(ctx, filter)
	if result.Err() != nil {
		// Se o erro for "mongo: no documents in result", você pode retornar um erro personalizado ou simplesmente nil
		if result.Err() == mongo.ErrNoDocuments {
			return nil, nil
		}

		log.Println("Erro ao buscar o usuário pelo id:", result.Err())
		return nil, result.Err()
	}

	err = result.Decode(&user)
	if err != nil {
		log.Println("Erro ao decodificar o usuário:", err)
		return nil, err
	}

	return &user, nil

}
func (r *UserRepository) Add(user pkgmodels.User) error {
	ctx := context.TODO()

	//inserir a data da criação da tarefa
	today := time.Now()
	eventRecord := pkgmodels.EventRecord{
		CreatedAt: &today,
	}
	user.EventRecord = &eventRecord

	_, err := r.Collection.InsertOne(ctx, user)
	if err != nil {
		log.Println("Erro ao inserir usuário:", err)
		return err
	}

	return nil
}
func (r *UserRepository) Update(id string, user pkgmodels.User) error {
	ctx := context.TODO()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Erro ao converter ID em ObjectID:", err)
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{
		"email":                 strings.ToLower(user.Email),
		"password":              user.Password,
		"eventRecord.updatedAt": time.Now()}}

	_, err = r.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Erro ao atualizar o usuário:", err)
		return err
	}

	return nil
}
func (r *UserRepository) Remove(id string) error {
	ctx := context.TODO()

	// Parseie a string do ID para um ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Erro ao converter ID em ObjectID:", err)
		return err
	}

	filter := bson.M{"_id": objectID}

	_, err = r.Collection.DeleteOne(ctx, filter)

	if err != nil {
		log.Println("Erro ao remover o usuário:", err)
		return err
	}
	return nil

}

func (r *UserRepository) FindByEmail(email string) (*pkgmodels.User, error) {
	var user pkgmodels.User

	ctx := context.TODO()

	filter := bson.M{"email": email}

	err := r.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Usuário não encontrado
			return nil, fmt.Errorf("usuário não encontrado")
		}
		log.Printf("Erro ao buscar usuário: %v", err)
		return nil, err
	}

	return &user, nil

}
