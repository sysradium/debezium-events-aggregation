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
	"github.com/sysradium/debezium-events-aggregation/service/internal/adapters/ephemeral"
	"github.com/sysradium/debezium-events-aggregation/service/internal/app"
	"github.com/sysradium/debezium-events-aggregation/service/internal/debezium/connectors"
	rm "github.com/sysradium/debezium-events-aggregation/service/internal/debezium/connectors/redis"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a := app.NewApplication(ephemeral.New())
	subClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})

	subscriber, err := redisstream.NewSubscriber(
		redisstream.SubscriberConfig{
			Client:        subClient,
			Unmarshaller:  rm.RedisStreamUnmarshaler{},
			ConsumerGroup: "test_consumer_group",
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		log.Fatal(err)
	}

	defer subscriber.Close()

	messages, err := subscriber.Subscribe(ctx, "debezium.monolith.auth_user")
	if err != nil {
		log.Panic(err)
	}

	handler := connectors.NewAuthUserHandler(messages, a)
	handler.Start(ctx)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

}
