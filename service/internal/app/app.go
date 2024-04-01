package app

import (
	"context"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/sysradium/debezium-events-aggregation/service/internal/app/commands"
	"github.com/sysradium/debezium-events-aggregation/service/internal/domain/users"
)

type App struct {
	Commands Commands
	UserRepo users.Repository
}

type Commands struct {
	CreateUser commands.CreateUserHandler
	DeleteUser commands.DeleteUserHandler
}

func NewApplication(u users.Repository) *App {
	// TODO: remove, for debug only
	go func() {
		timer := time.NewTicker(time.Second)
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				users, err := u.List(context.Background())
				if err != nil {
					fmt.Println("fuck", err)
				} else {
					spew.Dump(users)
				}
			}
		}
	}()
	return &App{
		Commands: Commands{
			CreateUser: commands.NewCreateUserHandler(u),
			DeleteUser: commands.NewDeleteUserHandler(u),
		},
		UserRepo: u,
	}
}
