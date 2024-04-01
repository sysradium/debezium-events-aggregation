package connectors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/sysradium/debezium-events-aggregation/service/internal/app"
	"github.com/sysradium/debezium-events-aggregation/service/internal/app/commands"
	"github.com/sysradium/debezium-events-aggregation/service/internal/debezium"
)

type AuthUserHandler struct {
	ch  <-chan *message.Message
	app *app.App
}

func NewAuthUserHandler(ch <-chan *message.Message, a *app.App) *AuthUserHandler {
	return &AuthUserHandler{
		ch:  ch,
		app: a,
	}
}

func (a *AuthUserHandler) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-a.ch:
				if !ok {
					return
				}

				if bytes.Equal(msg.Payload, []byte("default")) {
					fmt.Println("null payload, skip")
					msg.Ack()
					continue
				}
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
		return fmt.Errorf("could not umarshal payload: %w", err)
	}

	var au *User
	if err := json.Unmarshal(e.Payload.After, &au); err != nil {
		return fmt.Errorf("unable to unamrshal user: %w", err)
	}

	var bu *User
	if err := json.Unmarshal(e.Payload.Before, &bu); err != nil {
		return fmt.Errorf("unable to unamrshal user: %w", err)
	}

	switch e.Payload.Op {
	case debezium.OPERATION_SNAPSHOT, debezium.OPERATION_CREATE, debezium.OPERATION_UPDATE:
		_, err := a.app.Commands.CreateUser.Handle(
			context.Background(),
			commands.CreateUser{
				ID:          int64(au.ID),
				Password:    au.Password,
				Username:    au.Username,
				FirstName:   au.FirstName,
				LastName:    au.LastName,
				Email:       au.Email,
				IsSuperuser: au.IsSuperuser == 1,
			},
		)
		if err != nil {
			return err
		}
	case debezium.OPERATION_DELETE:
		_, err := a.app.Commands.DeleteUser.Handle(
			context.Background(),
			commands.DeleteUser{
				ID: int64(bu.ID),
			},
		)
		if err != nil {
			return err
		}
	}

	msg.Ack()

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
