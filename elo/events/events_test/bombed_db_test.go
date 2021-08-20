package events

import (
	"context"
	"testing"
)

func TestAllBombingsByDB(t *testing.T) {
	var dbcount int
	filecount := countAllBombings()

	row := store.Db.QueryRow(context.Background(), "select count(scoreaction.id) from scoreaction "+
		"left join players on actor=players.id "+
		"where actiontype='bombing';")
	err := row.Scan(&dbcount)
	if err != nil {
		panic(err.Error())
	}

	if filecount != dbcount {
		t.Errorf("%s failed: filecount %d != dbcount %d", t.Name(), filecount, dbcount)
	}
}
