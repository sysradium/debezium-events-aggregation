package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sysradium/debezium-events-aggregation/service/models"
)

func waitForConditionWithTimeout(condition func() bool, timeout time.Duration) bool {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	timeoutChan := time.After(timeout)
	for {
		select {
		case <-ticker.C:
			if condition() {
				return true
			}
		case <-timeoutChan:
			return false
		}
	}
}

type RedisStreamUnmarshaler struct{}

func (RedisStreamUnmarshaler) Unmarshal(values map[string]interface{}) (*message.Message, error) {
	payload := values["value"].(string)
	return  message.NewMessage(uuid.NewString(), []byte(payload)), nil
	// fmt.Println(values["key"].(string))
}

type TransactionManager struct {
	rClient *redis.Client
}

func (t *TransactionManager) Process(messages <-chan *message.Message) {

	for msg := range messages {
		var e models.DebeziumEvent
		if err := json.Unmarshal(msg.Payload, &e); err != nil {
			fmt.Println("error: ", err)
		}
		log.Printf(
			"---- received message: %v\n",
			e.Payload.Source.Table,
		)
		if e.Payload.Transaction.ID == "" {
			log.Println("no transaction information!!!")
			msg.Ack()
			continue
		}

		log.Printf("transaction: %v\n", spew.Sdump(e.Payload.Transaction))
		if err := t.rClient.SAdd(context.Background(), fmt.Sprintf("transaction_%v", e.Payload.Transaction.ID), string(msg.Payload)).Err(); err != nil {
			fmt.Printf("could not store transaction event: %v", err)
			continue
		}

		log.Printf("Payload: %v\n", string(e.Payload.After))
		msg.Ack()
	}
}

func (t *TransactionManager) ProcessTransacations(ctx context.Context, messages <-chan *message.Message) {

	p := func(msg *message.Message) error {
		var e models.TransactionMetadataEvent
		if err := json.Unmarshal(msg.Payload, &e); err != nil {
			return err
		}

		transactionKey := fmt.Sprintf("transaction_%v", e.Payload.ID)
		log.Printf(
			"---- received transaction: %v %v\n",
			transactionKey, e.Payload.Status,
		)

		if e.Payload.Status != "END" {
			log.Println("not an END transaction")
			return nil
		}

		log.Printf("WAITING FOR EVENTS TO COMMIT TRANSACTION: %v\n", e.Payload.ID)
		if !waitForConditionWithTimeout(
			func() bool {
				card, err := t.rClient.SCard(ctx, transactionKey).Result()
				if err != nil {
					return false
				}
				return card == e.Payload.EventCount
			},
			5*time.Second,
		) {
			log.Printf("CONDITIONS NOT MET, ABORTING: %v\n", e.Payload.ID)
			return errors.New("aborting")
		}

		log.Printf("COMMITING TRANSACTION: %v\n", e.Payload.ID)
		events, err := t.rClient.SMembers(ctx, transactionKey).Result()
		if err != nil {
			return err
		}

		log.Printf("COMMITING EVENTS: %v", spew.Sdump(events))
		t.rClient.Del(ctx, transactionKey)

		return nil

	}

	for msg := range messages {
		if err := p(msg); err != nil {
			log.Printf("error occured: %v, not acking", err)
			continue
		}
		msg.Ack()

	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	subClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})
	t := TransactionManager{
		rClient: subClient,
	}

	subscriber, err := redisstream.NewSubscriber(
		redisstream.SubscriberConfig{
			Client:        subClient,
			Unmarshaller:  RedisStreamUnmarshaler{},
			ConsumerGroup: "test_consumer_group",
		},
		watermill.NewStdLogger(false, false),
	)

	for _, topic := range []string{
		"debezium.applications.users",
		"debezium.applications.applications",
	} {
		messages, err := subscriber.Subscribe(ctx, topic)
		if err != nil {
			log.Panic(err)
		}

		go t.Process(messages)
	}

	transactions, err := subscriber.Subscribe(ctx, "debezium.transaction")
	if err != nil {
		log.Panic(err)
	}
	go t.ProcessTransacations(ctx, transactions)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	subscriber.Close()
}
