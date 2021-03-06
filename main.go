package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

// the topic and broker address are initialized as constants
const (
	topic          = "message-log"
	broker1Address = "localhost:9092"
	broker2Address = "localhost:9094"
	broker3Address = "localhost:9095"
)

func main() {
	// create a new context
	ctx := context.Background()
	// produce messages in a new go routine, since
	// both the produce and consume functions are
	// blocking
	//fmt.Println("tetsing")
	go produce(ctx)
	consume(ctx)
}

func produce(ctx context.Context) {
	// to produce messages
	topic := "my-topic"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10*time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte("one!")},
		kafka.Message{Value: []byte("two!")},
		kafka.Message{Value: []byte("three!")},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

func consume(ctx context.Context) {
	// to consume messages
	topic := "my-topic"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetReadDeadline(time.Now().Add(10*time.Second))
	batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

	b := make([]byte, 10e3) // 10KB max per message
	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		fmt.Println(string(b[:n]))
	}

	if err := batch.Close(); err != nil {
		log.Fatal("failed to close batch:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}
}

//func consume(ctx context.Context) {
//	// initialize a new reader with the brokers and topic
//	// the groupID identifies the consumer and prevents
//	// it from receiving duplicate messages
//	r := kafka.NewReader(kafka.ReaderConfig{
//		Brokers: []string{broker1Address, broker2Address, broker3Address},
//		Topic:   topic,
//		GroupID: "my-group",
//	})
//	for {
//		// the `ReadMessage` method blocks until we receive the next event
//		msg, err := r.ReadMessage(ctx)
//		if err != nil {
//			panic("could not read message " + err.Error())
//		}
//		// after receiving the message, log its value
//		fmt.Println("received: ", string(msg.Value))
//	}
//}

//func produce(ctx context.Context) {
//	// initialize a counter
//	i := 0
//
//	// intialize the writer with the broker addresses, and the topic
//	w := kafka.NewWriter(kafka.WriterConfig{
//		Brokers: []string{broker1Address, broker2Address, broker3Address},
//		Topic:   topic,
//	})
//
//	for {
//		// each kafka message has a key and value. The key is used
//		// to decide which partition (and consequently, which broker)
//		// the message gets published on
//		err := w.WriteMessages(ctx, kafka.Message{
//			Key: []byte(strconv.Itoa(i)),
//			// create an arbitrary message payload for the value
//			Value: []byte("this is message" + strconv.Itoa(i)),
//		})
//		if err != nil {
//			panic("could not write message " + err.Error())
//		}
//
//		// log a confirmation once the message is written
//		fmt.Println("writes:", i)
//		i++
//		// sleep for a second
//		time.Sleep(time.Second)
//	}
//}

