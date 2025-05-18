package main

import (
	"encoding/json" //helps in converting go structs to jason and vice versa
	"fmt"
	"log" //for logging error or messages

	"github.com/IBM/sarama"       //kafka client library used for sending and receiving messages to/from kafka
	"github.com/gofiber/fiber/v2" //web framework for creating REST API's , building web server
)

type Comment struct {
	Text string `form:"text" json:"text"`
}

func main() {
	app := fiber.New()                   //it creates a new fiber web server it is like starting our website backened
	api := app.Group("/api/v1")          //this creates a group of routes that all start with api/v1
	api.Post("/comments", createComment) //post endpoint
	app.Listen(":3000")
}

func ConnectProducer(brokerUrl []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()                           //creates a new kafka config
	config.Producer.Return.Successes = true                //after sending a message it confirms that it was successful
	config.Producer.RequiredAcks = sarama.WaitForAll       //wait for full confirmation from all the kafka brokers
	config.Producer.Retry.Max = 5                          //after sending messages if it fails it will try 5 times
	conn, err := sarama.NewSyncProducer(brokerUrl, config) //it creates a synchronous kafka producer
	if err != nil {
		return nil, err
	}
	return conn, nil //if theres an error, return it otherwise return kafka connection
}
func PushCommentToQueue(topic string, message []byte) error { //send messsages to kafka topic
	brokersUrl := []string{"localhost:29092"} //this port in localhost is used in docker compose
	producer, err := ConnectProducer(brokersUrl)
	if err != nil {
		return err
	}
	defer producer.Close() //closes the producer after this function finishes
	msg := &sarama.ProducerMessage{
		Topic: topic, //creates a kafka message with the topic and the message value
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err //sending a message to kafka where partition and offsett tell where in kafka the message was stored
	}
	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}

func createComment(c *fiber.Ctx) error { //We create a new Comment struct to store the user's data.
	cmt := new(Comment)
	if err := c.BodyParser(cmt); err != nil { //Reads the incoming JSON and fill the cmt struct
		log.Println(err)
		c.Status(400).JSON(&fiber.Map{ //If parsing fails, respond with a 400 Bad Request and show the error.
			"success": false,
			"message": err.Error(),
		})
		return err
	}
	cmtInBytes, err := json.Marshal(cmt) //convert the comment struct into JSON so we can send it to kafka
	if err != nil {
		log.Println(err)
		return err
	}
	PushCommentToQueue("comments", cmtInBytes) //Send the comment to the Kafka topic named comments

	err = c.JSON(&fiber.Map{ //Send a JSON response back to the client confirming it worked.

		"success": true,
		"message": "Comment pushed successfully",
		"comment": cmt,
	})
	if err != nil { //If sending the response fails, return a 500 Internal Server Error.
		c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Error creating product",
		})
		return err
	}
	return nil //All done! Everything worked, so return nil (no error).
}
