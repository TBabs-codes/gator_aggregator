package main

import (
	"context"
	"fmt"

	"github.com/TBabs-codes/gator_aggregator/internal/config"
	"github.com/TBabs-codes/gator_aggregator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmd_funcs map[string]func(*state, command) error
}

// Register command to be usable
func (c *commands) register(name string, f func(*state, command) error) {
	c.cmd_funcs[name] = f
	return
}

// Runs command
func (c *commands) run(s *state, cmd command) error {
	if f, ok := c.cmd_funcs[cmd.name]; ok {
		err := f(s, cmd)
		if err != nil {
			return err
		}
		return nil
	} else {
		return fmt.Errorf("command not found.")
	}
}

// Condenses user loggedIn checked for all commands that require user.
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {

	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUser)
		if err != nil {
			return fmt.Errorf("Invalid user logged in. Login again.")
		}

		err = handler(s, cmd, user)
		if err != nil {
			return err
		}
		return nil
	}
}
