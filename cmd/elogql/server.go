package main

import (
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
	"github.com/wlbr/cselo/elo/aggregators/graph"
	"github.com/wlbr/cselo/elo/aggregators/graph/generated"
	"github.com/wlbr/cselo/elo/sources/postgresql"
)

const defaultPort = "8080"

var (
	//Version is a linker injected variable for a git revision info used as version info
	Version = "Unknown build"
	//BuildTimestamp is a linker injected variable for a buildtime timestamp used in version info
	BuildTimestamp = "unknown build timestamp."

	config *elo.Config
)

func main() {

	config = new(elo.Config)
	config.Initialize(Version, BuildTimestamp)
	defer config.CleanUp()

	if s, e := postgresql.NewPostgres(config); e != nil {
		log.Error("Error opening postgres connection: %v", e)
	} else {

		port := os.Getenv("PORT")
		if port == "" {
			port = defaultPort
		}

		srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{Source: s}}))

		//http.Handle("/", playground.Handler("GraphQL playground", "/query"))
		http.Handle("/query", srv)

		log.Warn("Starting up")
		//log.Warn("connect to http://localhost:%s/ for GraphQL playground", port)
		log.Fatal("%s", http.ListenAndServe("localhost:"+port, nil))

	}
}
