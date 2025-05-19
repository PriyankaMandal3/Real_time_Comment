package main

import (
	"fmt"
	"os"        //used for interacting with the oprerating system
	"os/signal" //captures the os signals
	"syscall"   //used to define specific os signals

	"github.com/IBM/sarama"
)

func main() {
	topic := "comments"                                         //sets the kafka topic name to listen to comments
	worker, err := connectConsumer([]string{"localhost:29092"}) //Calls the connectConsumer function to connect to Kafka running at localhost:29092.
	if err != nil {
		panic(err)
	}
	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest) //Creates a new consumer for the specified topic and partition. It starts consuming from the oldest message available.
	if err != nil {
		panic(err)
	}
	fmt.Println("consumer started")
	sigchan := make(chan os.Signal, 1)                      //prepare to listen for OS signals
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM) //registers the signals we want to listen for
	msgCount := 0                                           //Tracks how many messages received.
	doneCh := make(chan struct{})                           //used to stop the program safely                      //Used to stop the program safely.
	go func() {
		for {
			select {
			case err := <-consumer.Errors(): //Prints any errors from Kafka.
				fmt.Println(err)
			case msg := <-consumer.Messages():
				msgCount++                                                                                                             //Increases message count.
				fmt.Printf("Received message Count: %d: | Topic (%s) | Message(%s)\n", msgCount, string(msg.Topic), string(msg.Value)) //Prints the topic and the message received.
			case <-sigchan:
				fmt.Println("Interrruption detection")
				doneCh <- struct{}{} //When user interrupts, it prints a message and sends a signal to stop the goroutine.
				return
			}
		}
	}()
	<-doneCh                                      //Waits until the signal is received to shut down the program.
	fmt.Println("Processed", msgCount, "message") //Shows how many messages were processed before stopping.
	if err := worker.Close(); err != nil {        //Closes the Kafka consumer connection safely.
		panic(err)
	}
}
func connectConsumer(brokersUrl []string) (sarama.Consumer, error) { //Creates a new Kafka consumer connection.
	config := sarama.NewConfig() //Tells kafka to return errors using config.Consumer.Return.Errors = true.
	config.Consumer.Return.Errors = true
	conn, err := sarama.NewConsumer(brokersUrl, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
