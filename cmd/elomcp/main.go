package main

import (
	"context"
	"crypto/subtle"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
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
	server := newMCPServer(config.GitVersion, source)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	runHTTPServer(ctx, server, config.MCP.Port, config.MCP.Path, config.MCP.BasicAuthUser, config.MCP.BasicAuthPassword)
}

func newMCPServer(version string, source statisticsSource) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "cselo-stats",
		Version: version,
	}, nil)
	registerTools(server, source)
	return server
}

func runHTTPServer(ctx context.Context, server *mcp.Server, port, path, basicAuthUser, basicAuthPassword string) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if (basicAuthUser == "") != (basicAuthPassword == "") {
		log.Fatal("MCP BasicAuthUser and BasicAuthPassword must either both be set or both be empty")
	}

	mux := http.NewServeMux()
	var handler http.Handler = mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{Stateless: true})
	if basicAuthUser != "" {
		handler = basicAuth(handler, basicAuthUser, basicAuthPassword)
	}
	mux.Handle(path, handler)
	mux.HandleFunc("/healthz", func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusNoContent)
	})

	httpServer := &http.Server{
		Addr:              ":" + port,
		Handler:           elo.LogHTTPRequests("elomcp", mux),
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		<-ctx.Done()
		shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownContext); err != nil {
			log.Error("MCP HTTP server shutdown failed: %v", err)
		}
	}()

	log.Warn("Starting MCP HTTP server at http://%s%s", httpServer.Addr, path)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("MCP HTTP server failed: %v", err)
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
