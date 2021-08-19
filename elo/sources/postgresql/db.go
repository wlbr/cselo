package postgresql

import (
	"context"
	"fmt"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"

	"github.com/jackc/pgx/v4"
)

type Postgres struct {
	config       *elo.Config
	Db           *pgx.Conn
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
		s.Db, err = pgx.Connect(context.Background(), dbinfo)
		if err != nil {
			log.Error("Cannot open PostgresQL database: %v", err)
		}
		cfg.AddCleanUpFn(func() error {
			log.Info("Cleanup - closing PostgreSQL database connection")
			return s.Db.Close(context.Background())
		})
		log.Info("Established PostgreSQL database connection")
	} else {
		log.Warn("Not creating PostgreSQL source due to missing connection data")
		return nil, fmt.Errorf("Insufficient PostgreSQL connection data (%s)", err)
	}
	return s, err
}
