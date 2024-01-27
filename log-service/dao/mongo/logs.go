package mongo

import (
	"context"
	"log"
	"time"

	"log-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName         = "logs"
	collectionName = "logs"
)

type LogRepoImpl struct {
	client *mongo.Client
}

func NewLogRepo(client *mongo.Client) *LogRepoImpl {
	return &LogRepoImpl{
		client: client,
	}
}

func (l *LogRepoImpl) Insert(entry models.LogEntry) error {
	collection := l.client.Database(dbName).Collection(collectionName)
	_, err := collection.InsertOne(context.TODO(), models.LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting into", collectionName, err)
		return err
	}
	return nil
}

func (l *LogRepoImpl) All() ([]*models.LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := l.client.Database(dbName).Collection(collectionName)

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Error finding all docs:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*models.LogEntry
	for cursor.Next(ctx) {
		var entry models.LogEntry
		if err = cursor.Decode(&entry); err != nil {
			log.Println("Error reading log entry:", err)
			return nil, err
		}
		logs = append(logs, &entry)
	}
	return logs, nil
}

func (l *LogRepoImpl) GetOne(id string) (*models.LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := l.client.Database(dbName).Collection(collectionName)

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry models.LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}
