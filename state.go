package main

import (
	"database/sql"

	"github.com/zoumas/gator/internal/config"
	"github.com/zoumas/gator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func newState() (*state, error) {
	cfg, err := config.Read()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		return nil, err
	}

	return &state{
		cfg: cfg,
		db:  database.New(db),
	}, nil
}
