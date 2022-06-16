package server

import (
	"github.com/pararti/forump/cmd/store"
)

type serverForum struct {
	store *store.CommonStore
}

func NewServer() *serverForum {
	return &serverForum{
		store: store.New(),
	}
}
