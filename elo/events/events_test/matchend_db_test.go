package events

import (
	"context"
	"testing"
)

func TestMatchEndsByDB(t *testing.T) {
	var dbmatchescount int
	matchescount := countMatchEnds()

	row := store.Db.QueryRow(context.Background(), "select count(matches.id) from matches where completed ")
	err := row.Scan(&dbmatchescount)
	if err != nil {
		panic(err.Error())
	}

	//if matchescount != dbmatchescount { // can't be equal, because we cannot correctly count started, but not finished maps using grep only
	if matchescount != dbmatchescount {
		t.Errorf("%s failed: filecount %d != dbcount %d", t.Name(), matchescount, dbmatchescount)
	}
}
