package events

import (
	"context"
	"testing"
)

func TestHostageRescuesPerPlayerByDB(t *testing.T) {
	var dbcount int
	filecount := countHostageRescuePerPlayer(player)

	row := store.Db.QueryRow(context.Background(), "select count(scoreaction.id) from scoreaction "+
		"left join players on actor=players.id "+
		"where actiontype='rescue' and players.initialname=$1;", player)
	err := row.Scan(&dbcount)
	if err != nil {
		panic(err.Error())
	}

	if filecount != dbcount {
		t.Errorf("%s failed: filecount %d != dbcount %d", t.Name(), filecount, dbcount)
	}
}

func TestAllHostageRescuesByDB(t *testing.T) {
	var dbcount int

	filecount := countAllHostageRescues()

	row := store.Db.QueryRow(context.Background(), "select count(scoreaction.id) from scoreaction "+
		"left join players on actor=players.id "+
		"where actiontype='rescue';")
	err := row.Scan(&dbcount)
	if err != nil {
		panic(err.Error())
	}

	if filecount != dbcount {
		t.Errorf("%s failed: filecount %d != dbcount %d", t.Name(), filecount, dbcount)
	}
}
