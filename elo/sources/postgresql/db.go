package postgresql

import (
	"context"
	"fmt"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	config       *elo.Config
	Db           *pgxpool.Pool
	PlayersCache elo.PlayersCache
}

func NewPostgres(cfg *elo.Config) (*Postgres, error) {
	log.Info("Creating new PostgreSQL source")
	var err error
	s := &Postgres{config: cfg, PlayersCache: make(elo.PlayersCache)}

	//pg connectstring "postgres://user:password@host:port5432/dbname"
	dbinfo := "postgres://"
	if cfg.PostgreSQL.User != "" {
		dbinfo += cfg.PostgreSQL.User
		if cfg.PostgreSQL.Password != "" {
			dbinfo += ":" + cfg.PostgreSQL.Password
		}
		dbinfo += "@"
	}
	if cfg.PostgreSQL.Host != "" {
		dbinfo += cfg.PostgreSQL.Host
		if cfg.PostgreSQL.Port != "" {
			dbinfo += ":" + cfg.PostgreSQL.Port
		}
	} else {
		log.Warn("No PostgresQL host given")
		err = fmt.Errorf("No PostgresQL host given")
	}
	if cfg.PostgreSQL.Database != "" {
		dbinfo += "/" + cfg.PostgreSQL.Database
	} else {
		log.Warn("No PostgresQL database name given")
		err = fmt.Errorf("No PostgresQL database name given")
	}
	if err == nil {
		s.Db, err = pgxpool.Connect(context.Background(), dbinfo)
		if err != nil {
			log.Error("Cannot open PostgresQL database: %v", err)
			return nil, err
		}
		cfg.AddCleanUpFn(func() error {
			log.Info("Cleanup - closing PostgreSQL connection pool")
			s.Db.Close()
			return nil
		})
		log.Info("Established PostgreSQL connection pool")
	} else {
		log.Warn("Not creating PostgreSQL source due to missing connection data")
		return nil, fmt.Errorf("Insufficient PostgreSQL connection data (%s)", err)
	}
	return s, err
}
