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

	l := log.New(os.Stdout, "[main] ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	db, err := ephemeral.New()
	if err != nil {
		l.Fatal(err)
	}

	a := app.NewApplication(db)
	subClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})

	subscriber, err := redisstream.NewSubscriber(
		redisstream.SubscriberConfig{
			Client:         subClient,
			Unmarshaller:   rm.RedisStreamUnmarshaler{},
			ConsumerGroup:  "",
			FanOutOldestId: "0",
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		l.Fatal(err)
	}

	defer subscriber.Close()

	messages, err := subscriber.Subscribe(ctx, "debezium.monolith.auth_user")
	if err != nil {
		l.Panic(err)
	}

	handler := connectors.NewAuthUserHandler(messages, a)
	handler.Start(ctx)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	l.Print("started service, waiting for messages")
	<-sigs

}
