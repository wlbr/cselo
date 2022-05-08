package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
	"github.com/wlbr/cselo/elo/aggregators/graph/generated"
)

func (r *matchResolver) Server(ctx context.Context, obj *elo.Match) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *matchResolver) Duration(ctx context.Context, obj *elo.Match) (*time.Time, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Players(ctx context.Context) ([]*elo.Player, error) {
	c := r.Source.GetAllPlayers()
	var players []*elo.Player
	for _, p := range c {
		players = append(players, p)
	}
	return players, nil
}

func (r *queryResolver) Player(ctx context.Context, id string) (*elo.Player, error) {
	intid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Error("Cannot convert ID %v to int: %s", id, err)
	}
	p := r.Source.GetPlayerByID(intid)
	return p, nil
}

func (r *queryResolver) Matches(ctx context.Context) ([]*elo.Match, error) {
	panic(fmt.Errorf("not implemented"))
}

// Match returns generated.MatchResolver implementation.
func (r *Resolver) Match() generated.MatchResolver { return &matchResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type matchResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
