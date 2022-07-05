package server

import (
	"github.com/pararti/forump/cmd/store/postgres"
	"github.com/pararti/forump/internals/entity"
)

type ServerForum struct {
	store *store.DataBase
}

func NewServer(config *entity.PSQLConfig) (*ServerForum, error) {
	db, err := store.NewDB(config)
	if err != nil {
		return &ServerForum{}, err
	}
	return &ServerForum{
		store: db,
	}, nil
}
