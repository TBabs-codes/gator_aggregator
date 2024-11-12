package main

import (
	"context"
	"fmt"
	"time"

	"github.com/TBabs-codes/gator_aggregator/internal/database"
	"github.com/google/uuid"
)


//Logs in user
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) > 1 {
		return fmt.Errorf("too many agrguments, only include username with login command.")
	} else if len(cmd.args) == 0 {
		return fmt.Errorf("no username provided.")
	}

	//Checks if user exists in DB
	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	//Sets the current user in the config file
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("User: %v, has logged in.\n", cmd.args[0])

	return nil
}

//Registers user and logs user in. username must be unique.
func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no name was provided for registration.")
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("User has been created. Username: %v\n", user.Name)

	return nil
}

//Displays all users
func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	if len(users) == 0 {
		fmt.Println("No users found in DB.")
		return nil
	}

	for i, user := range users {
		users[i] = "* " + users[i]
		if user == s.cfg.CurrentUser {
			users[i] += " (current)"
		}
	}

	for _, user := range users {
		fmt.Println(user)
	}
	return nil
}

//Deletes all users, feeds and posts from database.
func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't delete users: %w", err)
	}
	fmt.Println("Database reset successfully!")
	return nil
}

