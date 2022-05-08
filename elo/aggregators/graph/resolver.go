package graph

//go:generate go run github.com/99designs/gqlgen generate

import "github.com/wlbr/cselo/elo/sources/postgresql"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Source *postgresql.Postgres
}
