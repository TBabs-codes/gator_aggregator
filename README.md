Gator Aggregator helps users stay up to date on the favorite RSS feeds. See all the commands below to learn how to use it.

Commands:

Login:
    input: "go run . <login> <username>"

    description: logs in user using username


Register:

    input: go run . <register> <unique_name>

    description: creates a new user using provided name.


cmds.register("reset", handlerReset)
cmds.register("users", handlerUsers)
cmds.register("agg", handlerAgg)
cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
cmds.register("feeds", handlerFeeds)
cmds.register("follow", middlewareLoggedIn(handlerFollow))
cmds.register("following", middlewareLoggedIn(handlerFollowing))
cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
cmds.register("browse", middlewareLoggedIn(handlerBrowse))


Requirements:
Postgres and GO

Install instructions using go install:

Config file setup and how to run program.




