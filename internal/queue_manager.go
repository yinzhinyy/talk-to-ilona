package internal

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/mongodb/mongo-go-driver/bson"
)

var (
	topic              = "talk-to-ilona"
	group              = "mongo"
	bootstrapServers   = GetEnv("BOOTSTRAP_SERVERS")
	producer           *kafka.Producer
	producerInitialize bool
	consumer           *kafka.Consumer
	lock               sync.Mutex
	db                 = LoadDB()
)

// StartConsume creates consumer, pulls data from kafka and saves to database
// should run in a standalone process
func StartConsume() {
	sigchan := make(chan os.Signal, 1)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"group.id":          group,
		// "debug":             "all",
	})
	if err != nil {
		log.Fatal("failed to create kafka consumer")
		os.Exit(1)
	}

	err = c.SubscribeTopics(append([]string(nil), topic), nil)

	run := true
	for run == true {
		select {
		case sig := <-sigchan:
			log.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				log.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					log.Printf("%% Headers: %v\n", e.Headers)
				}
				var m *MessageBody
				err := json.Unmarshal(e.Value, &m)
				if err != nil {
					log.Fatal("json decode error")
					continue
				}
				saveToDB(*m)
			case kafka.Error:
				// Errors should generally be considered as informational, the client will try to automatically recover
				log.Fatalf("Error: %v\n", e)
			default:
				log.Printf("Ignored %v\n", e)
			}
		}
	}

	log.Printf("Closing consumer\n")
	c.Close()
}

func saveToDB(message MessageBody) {
	userName := message.UserName
	createTime := int64(message.CreateTime)
	createDate := time.Unix(createTime, 0)
	filter := bson.M{
		DBColumnUserName: userName,
		DBColumnYear:     createDate.Year(),
		DBColumnMonth:    createDate.Month(),
		DBColumnDay:      createDate.Day(),
	}
	var record bson.M
	record = db.Find(DBName, DBTableActivity, filter)
	activity := bson.M{
		DBColumnCreateTime: message.CreateTime,
		DBColumnActivity:   message.Activity,
	}
	if record == nil {
		record = bson.M{
			DBColumnUserName:   userName,
			DBColumnYear:       createDate.Year(),
			DBColumnMonth:      createDate.Month(),
			DBColumnDay:        createDate.Day(),
			DBColumnActivities: bson.A{activity},
		}
		db.Save(DBName, DBTableActivity, record)
	} else {
		activities := append(record[DBColumnActivities].(bson.A), activity)
		db.Update(DBName, DBTableActivity, filter, bson.M{
			"$set": bson.M{
				DBColumnActivities: activities,
			},
		})
	}
}

// Produce creates producer if needed and push data to kafka
func Produce(message *MessageBody) {
	initProducer()
	bytes, err := json.Marshal(message)
	if err != nil {
		log.Fatal("json encoding error")
		return
	}

	producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          bytes,
		Key:            []byte(message.UserName),
	}, nil)

}

func initProducer() {
	if !producerInitialize {
		lock.Lock()
		defer lock.Unlock()
		if !producerInitialize {
			p, err := kafka.NewProducer(&kafka.ConfigMap{
				"bootstrap.servers": bootstrapServers,
				// "debug":             "all",
			})
			if err != nil {
				log.Fatal("failed to create kafka producer")
				os.Exit(1)
			}
			producer = p
			log.Print("kafka producer created")
		}
		producerInitialize = true
	}
}
