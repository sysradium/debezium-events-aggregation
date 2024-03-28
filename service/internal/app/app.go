package app

import (
	"github.com/sysradium/debezium-events-aggregation/service/internal/app/commands"
	"github.com/sysradium/debezium-events-aggregation/service/internal/domain/users"
)

type App struct {
	Commands Commands
}

type Commands struct {
	CreateUser commands.CreateUserHandler
}

func NewApplication(u users.Repository) *App {
	return &App{
		Commands: Commands{
			CreateUser: commands.NewCreateUserHandler(u),
		},
	}
}
