package internal

import (
	"context"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type MongoDB struct {
	URL    string
	Client *mongo.Client
}

func (mongo *MongoDB) Save(db string, table string, document interface{}) interface{} {
	c, err := getClient(mongo)
	if err != nil {
		log.Fatal(err)
		return 0
	}
	collection := c.Database(db).Collection(table)
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	res, err := collection.InsertOne(ctx, document)
	if err != nil {
		log.Fatal(err)
		return 0
	}
	return res.InsertedID
}

func (mongo *MongoDB) Update(db string, table string, filter interface{}, document interface{}) interface{} {
	c, err := getClient(mongo)
	if err != nil {
		log.Fatal(err)
		return 0
	}
	collection := c.Database(db).Collection(table)
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	res, err := collection.UpdateOne(ctx, filter, document)
	if err != nil {
		return 0
	}
	return res.UpsertedID
}

func (mongo *MongoDB) Find(db string, table string, filter interface{}) bson.M {
	var result bson.M
	c, err := getClient(mongo)
	if err != nil {
		log.Fatal(err)
		return result
	}
	collection := c.Database(db).Collection(table)
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	collection.FindOne(ctx, filter).Decode(&result)
	return result
}

func getClient(config *MongoDB) (*mongo.Client, error) {
	if config.Client == nil {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		client, err := mongo.Connect(ctx, config.URL)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		config.Client = client
	}
	return config.Client, nil
}
