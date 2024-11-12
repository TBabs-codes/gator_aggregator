package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/TBabs-codes/gator_aggregator/internal/config"
	"github.com/TBabs-codes/gator_aggregator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Thanks for jamming me with commands.")

	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	db, err := sql.Open("postgres", "postgres://thomasbabcock:@localhost:5432/gator_aggregator?sslmode=disable")
	dbQueries := database.New(db)

	s := state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		cmd_funcs: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) >= 2 {
		cmd := command{
			name: os.Args[1],
			args: os.Args[2:],
		}
		err := cmds.run(&s, cmd)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	} else {
		fmt.Println("No command given")
		os.Exit(1)
	}
}
