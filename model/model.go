package model

import (
	"github.com/go-pg/pg"
	_ "github.com/go-pg/pg/orm"
)

type Model struct {
	DB *pg.DB
}

func NewDB() *Model {
	db := pg.Connect(&pg.Options{
		User:     "demo",
		Password: "demo",
	})

	return &Model{db}
}
