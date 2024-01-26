package mongo

import (
	"context"
	"log"
	"time"
)

func (l *LogRepo) Insert(entry LogEntry) error {
	log.Println("Insert:", entry.Name, entry.Data)
	collection := l.client.Database(dbName).Collection(collectionName)
	_, err := collection.InsertOne(context.TODO(), LogEntry{
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
