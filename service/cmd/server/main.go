package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/redis/go-redis/v9"
	"github.com/sysradium/debezium-events-aggregation/service/internal/debezium/connectors"
	rm "github.com/sysradium/debezium-events-aggregation/service/internal/debezium/connectors/redis"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	subClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})
	t := connectors.New(subClient)

	subscriber, err := redisstream.NewSubscriber(
		redisstream.SubscriberConfig{
			Client:        subClient,
			Unmarshaller:  rm.RedisStreamUnmarshaler{},
			ConsumerGroup: "test_consumer_group",
		},
		watermill.NewStdLogger(false, false),
	)

	for _, topic := range []string{
		"debezium.auth_user",
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
