Gator Aggregator helps users stay up to date on the favorite RSS feeds. See all the commands below to learn how to use it.

Commands:

Login:
    input: "go run . login <username>"
    
    description: logs in user using username


Register:

    input: go run . register <unique_name>

    description: creates a new user using provided name.


Reset:

    input: go run . reset

    description: resets all the data in the database.

Users:

    input: go run . users

    description: shows all users in the database and specifies who is logged in.

Agg: (runs continuously until commanded to stop)

    input: go run . agg <update.timing> 
            update timing examples: 1h, 5m, 3s

    description: starts the aggregation process, parsing RSS feeds into individual posts and storing them in the database. 

AddFeed:

    input: go run . addfeed <name> <url>

    description: adds a feed with name and url specified. Current user will also automatically start follow feed if successfully added.

Feeds:

    input: go run . feeds

    description: prints all the feeds in the database.

Follow:

    input: go run . follow <url>

    description: current user starts following the feed associated with that url. If url is not associated to a feed you will need to add the feed with that url first.


Following:

    input: go run . following

    description: prints all the feeds the current user is following


Unfollow:

    input: go run . unfollow <url>

    desciption: user will no longer see updates from this feed when using browse.

Browse:

    input: go run . browse <number of posts>

    description: the most up to date posts will be printed option to specify how many posts the users would like to be printed.


Requirements:
Postgres and GO

Install instructions using go install:

Config file setup and how to run program.




