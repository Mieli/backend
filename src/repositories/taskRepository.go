package repositories

import (
	"context"
	"log"
	"os"
	"time"

	pkgmodels "github.com/mieli/backend/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepository struct {
	DB         *mongo.Client
	Collection *mongo.Collection
}

func NewTaskRepository(db *mongo.Client) *TaskRepository {
	return &TaskRepository{
		DB:         db,
		Collection: db.Database(os.Getenv("MONGO_DATABASE")).Collection("tasks"),
	}

}

func (r *TaskRepository) FindAll() ([]pkgmodels.Task, error) {
	// Criar uma lista vazia de tarefas
	var tasks []pkgmodels.Task

	ctx := context.TODO()
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Erro ao buscar tarefas:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var task pkgmodels.Task
		if err := cursor.Decode(&task); err != nil {
			log.Println("Erro ao decodificar tarefa:", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil

}

func (r *TaskRepository) FindById(id string) (*pkgmodels.Task, error) {
	var task pkgmodels.Task

	// Parseie a string do ID para um ObjectID
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

		log.Println("Erro ao buscar a tarefa pelo id:", result.Err())
		return nil, result.Err()
	}

	err = result.Decode(&task)
	if err != nil {
		log.Println("Erro ao decodificar a tarefa:", err)
		return nil, err
	}

	return &task, nil

}

func (r *TaskRepository) Add(task pkgmodels.Task) error {

	ctx := context.TODO()

	//inserir a data da criação da tarefa
	today := time.Now()
	eventRecord := pkgmodels.EventRecord{
		CreatedAt: &today,
	}
	task.EventRecord = &eventRecord

	_, err := r.Collection.InsertOne(ctx, task)
	if err != nil {
		log.Println("Erro ao inserir tarefa:", err)
		return err
	}

	return nil
}

func (r *TaskRepository) Update(id string, task pkgmodels.Task) error {

	ctx := context.TODO()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Erro ao converter ID em ObjectID:", err)
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"title": task.Title, "eventRecord.updatedAt": time.Now()}}

	_, err = r.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Erro ao atualizar a tarefa:", err)
		return err
	}

	return nil
}

func (r *TaskRepository) Remove(id string) error {
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
		log.Println("Erro ao remover a tarefa:", err)
		return err
	}
	return nil

}
