package main

import (
	"context"
	"crypto/subtle"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
	"github.com/wlbr/cselo/elo/aggregators/graph"
	"github.com/wlbr/cselo/elo/aggregators/graph/generated"
	"github.com/wlbr/cselo/elo/analytics"
	"github.com/wlbr/cselo/elo/sources/postgresql"
)

var (
	// Version is a linker injected variable for a git revision info used as version info.
	Version = "Unknown build"
	// BuildTimestamp is a linker injected variable for a buildtime timestamp used as version info.
	BuildTimestamp = "unknown build timestamp."
)

func main() {
	config := new(elo.Config)
	config.Initialize(Version, BuildTimestamp)
	defer config.CleanUp()

	log.Print("Starting up.  Version: %s.   Timestamp: %s    Working directory:%s ", config.GitVersion, config.BuildTimeStamp, config.WorkingDirectory)
	start := time.Now()
	defer func() {
		log.Warn("Shutting down, been up for %s", time.Since(start))
	}()

	source, err := postgresql.NewPostgres(config)
	if err != nil {
		log.Fatal("Error opening PostgreSQL connection: %v", err)
	}
	service := analytics.New(source)
	mux, err := newHandler(config.GitVersion, service, config.MCP.BasicAuthUser, config.MCP.BasicAuthPassword)
	if err != nil {
		log.Fatal("Invalid server configuration: %v", err)
	}
	if config.MCP.Path != "" && config.MCP.Path != "/mcp" {
		log.Fatal("MCP.Path must be /mcp for elointel")
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	runHTTPServer(ctx, mux, config.MCP.Host, config.MCP.Port)
}

func newMCPServer(version string, service *analytics.Service) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "cselo-stats",
		Version: version,
	}, nil)
	registerTools(server, service)
	return server
}

func newHandler(version string, service *analytics.Service, basicAuthUser, basicAuthPassword string) (http.Handler, error) {
	if (basicAuthUser == "") != (basicAuthPassword == "") {
		return nil, fmt.Errorf("MCP BasicAuthUser and BasicAuthPassword must either both be set or both be empty")
	}

	mux := http.NewServeMux()
	server := newMCPServer(version, service)
	var handler http.Handler = mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{Stateless: true})
	if basicAuthUser != "" {
		handler = basicAuth(handler, basicAuthUser, basicAuthPassword)
	}
	mux.Handle("/mcp", handler)
	mux.Handle("/graphql", graphqlHandler(service))
	return mux, nil
}

func graphqlHandler(service *analytics.Service) http.Handler {
	return handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &graph.Resolver{Analytics: service},
	}))
}

func runHTTPServer(ctx context.Context, mux http.Handler, host, port string) {
	if host == "" {
		host = "localhost"
	}
	httpServer := &http.Server{
		Addr:              net.JoinHostPort(host, port),
		Handler:           elo.LogHTTPRequests("elointel", mux),
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		<-ctx.Done()
		shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownContext); err != nil {
			log.Error("HTTP server shutdown failed: %v", err)
		}
	}()

	log.Warn("Starting HTTP server at http://%s (MCP: /mcp, GraphQL: /graphql)", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("HTTP server failed: %v", err)
	}
}

func basicAuth(next http.Handler, username, password string) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		providedUsername, providedPassword, ok := request.BasicAuth()
		if !ok ||
			subtle.ConstantTimeCompare([]byte(providedUsername), []byte(username)) != 1 ||
			subtle.ConstantTimeCompare([]byte(providedPassword), []byte(password)) != 1 {
			writer.Header().Set("WWW-Authenticate", `Basic realm="cselo-mcp"`)
			http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(writer, request)
	})
}

func validateLimit(limit int) (int, error) {
	if limit == 0 {
		return 20, nil
	}
	if limit < 1 || limit > 100 {
		return 0, fmt.Errorf("limit must be between 1 and 100")
	}
	return limit, nil
}

func validateOffset(offset int) error {
	if offset < 0 {
		return fmt.Errorf("offset must not be negative")
	}
	return nil
}
