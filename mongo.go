package talk

import (
	"context"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
)

const dbURL = "mongodb://35.247.89.169:27017"

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
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	collection := c.Database(db).Collection(table)
	res, err := collection.InsertOne(ctx, document)
	if err != nil {
		return 0
	}
	return res.InsertedID
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
