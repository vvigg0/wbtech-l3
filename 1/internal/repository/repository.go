package repository

import (
	"errors"

	"github.com/wb-go/wbf/dbpg"
)

var ErrNoNotification error = errors.New("такого уведомления нет")

type Repository struct {
	*dbpg.DB
}

func New(db *dbpg.DB) *Repository {
	return &Repository{db}
}
