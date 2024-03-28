package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/sysradium/debezium-events-aggregation/service/internal/debezium"
)

type AuthUserHandler struct {
	ch <-chan *message.Message
}

func NewAuthUserHandler(ch <-chan *message.Message) *AuthUserHandler {
	return &AuthUserHandler{
		ch: ch,
	}
}

func (a *AuthUserHandler) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-a.ch:
				if err := a.Handle(msg); err != nil {
					log.Printf("unable to process message: %v", err)
				}
			}
		}
	}()
}

func (a *AuthUserHandler) Handle(msg *message.Message) error {
	var e debezium.DebeziumEvent
	if err := json.Unmarshal(msg.Payload, &e); err != nil {
		return err
	}

	var u User
	if err := json.Unmarshal(e.Payload.After, &u); err != nil {
		return fmt.Errorf("unable to unamrshal user: %w", err)
	}

	switch e.Payload.Op {
	case debezium.OPERATION_CREATE:
		fmt.Println("creating new user")
	case debezium.OPERATION_SNAPSHOT:
		fmt.Println("a user from snapshot")
	case debezium.OPERATION_DELETE:
		fmt.Println("remove user")
	case debezium.OPERATION_UPDATE:
		fmt.Println("changing user")
	}

	return nil
}

type User struct {
	ID          int    `json:"id"`
	Password    string `json:"password"`
	LastLogin   *int64 `json:"last_login"`
	IsSuperuser int    `json:"is_superuser"`
	Username    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	IsStaff     uint   `json:"is_staff"`
	IsActive    uint   `json:"is_active"`
	DateJoined  int64  `json:"date_joined"`
}
