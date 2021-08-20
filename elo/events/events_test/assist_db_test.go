package events

import (
	"context"
	"testing"
)

func TestAssistsPerPlayerByDB(t *testing.T) {
	var dbcount int
	filecount := countAssistsPerPlayer(player)

	row := store.Db.QueryRow(context.Background(), "select count(assists.id) from assists "+
		"left join players on actor=players.id "+
		"where players.initialname=$1;", player)
	err := row.Scan(&dbcount)
	if err != nil {
		panic(err.Error())
	}

	if filecount != dbcount {
		t.Errorf("%s failed: filecount %d != dbcount %d", t.Name(), filecount, dbcount)
	}
}

func TestAllAssistsByDB(t *testing.T) {
	var dbcount int
	filecount := countAllAssists()

	row := store.Db.QueryRow(context.Background(), "select count(assists.id) from assists "+
		"left join players on actor=players.id;")
	err := row.Scan(&dbcount)
	if err != nil {
		panic(err.Error())
	}

	if filecount != dbcount {
		t.Errorf("%s failed: filecount %d != dbcount %d", t.Name(), filecount, dbcount)
	}
}
