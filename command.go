package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zoumas/gator/internal/database"
)

type command struct {
	name string
	args []string
}

type commands struct {
	m map[string]func(*state, command) error
}

func newCommands() *commands {
	m := map[string]func(*state, command) error{
		"login":    handlerLogin,
		"register": handlerRegister,
	}
	return &commands{
		m: m,
	}
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.m[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.m[cmd.name]
	if !ok {
		return fmt.Errorf("command %v does not exist", cmd.name)
	}
	return f(s, cmd)
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("expected username")
	}
	name := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("user %q does not exist: %v", name, err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return err
	}
	fmt.Println("logged in as", name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("expected username")
	}
	name := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err == nil {
		return fmt.Errorf("user %q already exists", name)
	}

	now := time.Now().UTC()

	_, err = s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
	})
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return err
	}
	fmt.Println("user", name, "created")
	return err
}