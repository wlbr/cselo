package main

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
	"github.com/wlbr/cselo/elo/aggregators/graph"
	"github.com/wlbr/cselo/elo/aggregators/graph/generated"
	"github.com/wlbr/cselo/elo/sources/postgresql"
)

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

	log.Print("Starting up.  Version: %s.   Timestamp: %s    Working directory:%s ", config.GitVersion, config.BuildTimeStamp, config.WorkingDirectory)
	start := time.Now()
	defer func() {
		log.Warn("Shutting down, been up for %s", time.Since(start))
	}()

	if s, e := postgresql.NewPostgres(config); e != nil {
		log.Error("Error opening postgres connection: %v", e)
	} else {

		srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{Source: s}}))

		//http.Handle("/", playground.Handler("GraphQL playground", "/query"))
		http.Handle("/query", elo.LogHTTPRequests("elogql", srv))

		log.Warn("Starting GraphQL server at http://localhost:%s/query", config.GraphQL.Port)
		//log.Warn("connect to http://localhost:%s/ for GraphQL playground", port)
		log.Fatal("%s", http.ListenAndServe("localhost:"+config.GraphQL.Port, nil))

	}
}
