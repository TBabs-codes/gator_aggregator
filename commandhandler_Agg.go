package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/TBabs-codes/gator_aggregator/internal/database"
	"github.com/google/uuid"
)

//Starts forever loop that parses RSS feeds until user interupts program.
func handlerAgg(s *state, cmd command) error {

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("Unable to convert string to time duration. Valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\".")
	}
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		fmt.Println("SCRAPING Feeds!")
		err = scrapeFeeds(s)
	}

}

//Finds next feed to be updated then pulls information offline and parses it into posts table.
func scrapeFeeds(s *state) error {

	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	err = s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return err
	}

	feed, err := fetchFeed(context.Background(), nextFeed.Url)

	// fmt.Println("Channel Title: %v", feed.Channel.Title)
	// fmt.Println("Channel Description: %v", feed.Channel.Description)

	// for i, item := range feed.Channel.Item {
	// 	fmt.Printf("Item %v: %v\n", i, item.Title)
	// 	fmt.Printf("	Description: %v\n", item.Description)
	// 	fmt.Printf("			URL: %v\n", item.Link)
	// }
	
	for _, item := range feed.Channel.Item {
		
		pub_data, err := time.Parse("Mon, 02 Jan 2006 15:04:05 +0000", item.PubDate)
		if err != nil {
			return fmt.Errorf("Couldn't convert published time to SQL TIMESTAMP")
		}

		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: pub_data,
			FeedID:      nextFeed.ID,
		})
		if err != nil {
			fmt.Println("Error creating post: ", err)
		}
		fmt.Println("Item added to posts successfully.")
	}

	return nil
}

//Displays most up-to-date posts from user's following.
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