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
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-a.ch:
			if !ok {
				return
			}

			if isTumbstoneEvent(msg) {
				msg.Ack()
				continue
			}

			if err := a.Handle(msg); err != nil {
				log.Printf("unable to process message: %v", err)
				msg.Nack()
				continue
			}

			msg.Ack()
		}
	}
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

	return nil
}

func isTumbstoneEvent(m *message.Message) bool {
	return bytes.Equal(m.Payload, []byte("default"))
}
