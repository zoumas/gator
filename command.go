package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zoumas/gator/internal/database"
	"github.com/zoumas/gator/internal/rss"
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
		"login":     handlerLogin,
		"register":  handlerRegister,
		"reset":     handleReset,
		"users":     handleUsers,
		"agg":       auth(handleAgg),
		"addfeed":   auth(handlerAddFeed),
		"feeds":     handlerFeeds,
		"follow":    auth(handlerFollow),
		"following": auth(handlerFollowing),
		"unfollow":  auth(handlerUnfollow),
		"scrape":    auth(scrapeFeed),
	}
	return &commands{
		m: m,
	}
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

func handleReset(s *state, _ command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not delete all users: %v", err)
	}
	fmt.Println("true reset: deleted all users")
	return nil
}

func handleUsers(s *state, _ command) error {
	us, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, u := range us {
		if u.Name == s.cfg.CurrentUserName {
			fmt.Println("*", u.Name, "(current)")
		} else {
			fmt.Println("*", u.Name)
		}
	}

	return nil
}

func handleAgg(s *state, cmd command, u database.User) error {
	if len(cmd.args) == 0 {
		return errors.New("expected duration between requests")
	}
	d, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	ticker := time.NewTicker(d)
	for ; ; <-ticker.C {
		scrapeFeed(s, cmd, u)
	}
}

func handlerAddFeed(s *state, cmd command, u database.User) error {
	if len(cmd.args) != 2 {
		return errors.New("expected name and url of rss feed")
	}
	name := cmd.args[0]
	url := cmd.args[1]

	now := time.Now().UTC()
	f, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
		Url:       url,
		UserID:    u.ID,
	})
	if err != nil {
		return err
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    u.ID,
		FeedID:    f.ID,
	})
	return err
}

func handlerFeeds(s *state, _ command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, f := range feeds {
		fmt.Println(f.Name, f.Url, f.Owner)
	}
	return nil
}

func handlerFollow(s *state, cmd command, u database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("expected url")
	}
	url := cmd.args[0]

	f, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	r, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    u.ID,
		FeedID:    f.ID,
	})
	if err != nil {
		return err
	}

	fmt.Println(r.UserName, "followed", r.FeedName)

	return nil
}

func handlerFollowing(s *state, _ command, u database.User) error {
	ffs, err := s.db.GetFeedFollowsForUser(context.Background(), u.ID)
	if err != nil {
		return err
	}

	for _, ff := range ffs {
		fmt.Println(ff.FeedName)
	}
	return nil
}

func auth(handler func(s *state, cmd command, u database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		u, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, u)
	}
}

func handlerUnfollow(s *state, cmd command, u database.User) error {
	if len(cmd.args) == 0 {
		return errors.New("expected url of feed to unfollow")
	}
	url := cmd.args[0]
	return s.db.UnfollowFeedForUser(context.Background(), database.UnfollowFeedForUserParams{
		UserID: u.ID,
		Url:    url,
	})
}

func scrapeFeed(s *state, cmd command, u database.User) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return err
	}

	rssFeed, err := rss.Fetch(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	fmt.Println(rssFeed.Channel.Title, rssFeed.Channel.Description)
	for _, item := range rssFeed.Channel.Item {
		fmt.Println(item.Title)
	}
	return nil
}
