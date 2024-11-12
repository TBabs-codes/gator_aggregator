package main

import (
	"context"
	"fmt"
	"time"

	"github.com/TBabs-codes/gator_aggregator/internal/database"
	"github.com/google/uuid"
)

// Adds a feed with a specified name and url. Additional the user who performs this action will follow that feed.
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
