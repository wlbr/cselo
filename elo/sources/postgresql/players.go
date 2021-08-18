package postgresql

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/wlbr/commons/log"
	"github.com/wlbr/cs-elo/elo"
)

func (s *Postgres) GetPlayerByID(id int64) *elo.Player {
	log.Info("Getting PLAYER by id from PostgreSQL database. id='%v'", id)

	p := &elo.Player{}
	row := s.Db.QueryRow(context.Background(), `SELECT id, initialname, steamid, profileid FROM players WHERE id=$1`,
		id)
	err := row.Scan(&p.ID, &p.Name, &p.SteamID, &p.ProfileID)
	switch {
	case err == pgx.ErrNoRows:
		log.Info("Cannot find PLAYER with id '%d' in PostgresQL database: %v", id, err)
	case err != nil:
		log.Error("Cannot read from PostgresQL database: %v", err)
	}
	return p
}

func (s *Postgres) GetPlayerBySteamID(steamid string) *elo.Player {
	log.Info("Getting PLAYER by steamid from PostgreSQL database. steamid='%v'", steamid)

	p := &elo.Player{}
	row := s.Db.QueryRow(context.Background(), `SELECT id, initialname, steamid, profileid FROM players WHERE steamid=$1`,
		steamid)
	err := row.Scan(&p.ID, &p.Name, &p.SteamID, &p.ProfileID)
	switch {
	case err == pgx.ErrNoRows:
		log.Info("Cannot find PLAYER with steamid '%s' in PostgresQL database: %v", steamid, err)
	case err != nil:
		log.Error("Cannot read from PostgresQL database: %v", err)
	}
	return p
}

func (s *Postgres) GetAllPlayers() elo.PlayersCache {
	log.Info("Getting all players.")

	rows, err := s.Db.Query(context.Background(), "SELECT id, initialname, steamid, profileid FROM players;")
	defer rows.Close()
	if err != nil {
		log.Error("Cannot read all players from PostgresQL database: %v", err)
	}

	for rows.Next() {
		p := &elo.Player{}
		err = rows.Scan(&p.ID, &p.Name, &p.SteamID, &p.ProfileID)
		switch {
		case err == pgx.ErrNoRows:
			log.Info("Cannot find players in PostgresQL database: %v", err)
		case err != nil:
			log.Error("Cannot read player from PostgresQL database: %v", err)
		default:
			if _, present := s.PlayersCache[int(p.ID)]; !present {
				s.PlayersCache[int(p.ID)] = p
			}
		}
	}

	return s.PlayersCache
}
