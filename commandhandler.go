package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/TBabs-codes/gator_aggregator/internal/config"
	"github.com/TBabs-codes/gator_aggregator/internal/database"
	"github.com/google/uuid"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

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

type commands struct {
	cmd_funcs map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmd_funcs[name] = f
	return
}

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

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't delete users: %w", err)
	}
	fmt.Println("Database reset successfully!")
	return nil
}

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

func handlerAddFeed(s *state, cmd command, user database.User) error {

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("feed creation failed")
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Unable to follow URL.")
	}

	fmt.Println("New Feed created: ", feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found in DB.")
		return nil
	}
	fmt.Println("Feeds found:")
	fmt.Println("=================================================================")
	for _, feed := range feeds {
		printFeed(feed)
		fmt.Println("=================================================================")
	}

	return nil
}

func printFeed(feed database.GetFeedsRow) {
	fmt.Printf("* Feed Name:     %s\n", feed.FeedName)
	fmt.Printf("* Created By:    %s\n", feed.CreatedBy)
	fmt.Printf("* URL:           %s\n", feed.URL)
}

func handlerFollow(s *state, cmd command, user database.User) error {
	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Unable to find feed with URL provided.")
	}

	feed_follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Unable to follow URL.")
	}

	fmt.Println(feed_follow.UserName, " now follows ", feed_follow.FeedName, "!!")

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	following, err := s.db.GetFeedFollowForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error retrieving user's following records")
	}
	if len(following) == 0 {
		fmt.Println("You are not following anyone.")
		return nil
	}

	fmt.Println("==========================")
	fmt.Println(" You are following:")
	for _, follow := range following {
		fmt.Printf(" * %v\n", follow.FeedName)
	}
	fmt.Println("==========================")

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Invalid URL. Not Found in DB.")
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("unfollow request was unsuccessful.")
	}

	fmt.Printf("You no longer follow %v.\n", feed.Name)
	return nil
}

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

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int32
	if len(cmd.args) == 0 {
		limit = 2
	} else {
		lim, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("Unable to convert string into integer.")
		}
		limit = int32(lim)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("Error retrieving posts from DB.")
	}

	fmt.Println("======================================================")
	for _, post := range posts {
		fmt.Printf("Title: %v\n", post.Title)
		fmt.Printf("Description: %v\n", post.Description)
		fmt.Printf("URL: %v\n", post.Url)
		fmt.Println("======================================================")
	}

	
	return nil
}
