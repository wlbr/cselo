package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo/aggregators/graph/generated"
	"github.com/wlbr/cselo/elo/aggregators/graph/model"
)

func (r *queryResolver) Players(ctx context.Context) ([]*model.Player, error) {
	c := r.Source.GetAllPlayers()
	var players []*model.Player
	for _, p := range c {
		players = append(players, &model.Player{ID: fmt.Sprint(p.ID), Name: p.Name, Steamid: p.SteamID, Profileid: p.ProfileID})
	}
	return players, nil
}

func (r *queryResolver) Player(ctx context.Context, id string) (*model.Player, error) {
	intid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Error("Cannot convert ID %v to int: %s", id, err)
	}
	p := r.Source.GetPlayerByID(intid)
	return &model.Player{ID: fmt.Sprint(p.ID), Name: p.Name, Steamid: p.SteamID, Profileid: p.ProfileID}, nil
}

func (r *queryResolver) Matches(ctx context.Context) ([]*model.Match, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
