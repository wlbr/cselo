package events

import (
	"context"
	"testing"
)

func TestKillsPerPlayerByDB(t *testing.T) {
	var dbkillcount int
	killcount := countKillsPerPlayer(player)

	row := store.Db.QueryRow(context.Background(), "select count(kills.id) from kills "+
		"left join players on actor=players.id "+
		"where players.initialname=$1;", player)
	err := row.Scan(&dbkillcount)
	if err != nil {
		panic(err.Error())
	}

	if killcount != dbkillcount {
		t.Errorf("%s failed: filecount %d != dbcount %d", t.Name(), killcount, dbkillcount)
	}
}

func TestAllKillsByDB(t *testing.T) {
	var dbkillcount int
	killcount := countAllKills()

	row := store.Db.QueryRow(context.Background(), "select count(kills.id) from kills "+
		"left join players on actor=players.id")
	err := row.Scan(&dbkillcount)
	if err != nil {
		panic(err.Error())
	}

	if killcount != dbkillcount {
		t.Errorf("%s failed: filecount %d != dbcount %d", t.Name(), killcount, dbkillcount)
	}
}
