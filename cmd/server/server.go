package server

import (
	"github.com/pararti/forump/cmd/store/postgres"
	"github.com/pararti/forump/internals/entity"
)

type ServerForum struct {
	store *store.DataBase
}

func NewServer(config *entity.PSQLConfig) *ServerForum {
	return &ServerForum{
		store: store.NewDB(config),
	}
}
