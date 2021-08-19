package events

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/commander-cli/cmd"
)

var playersdefuses int = -1

func countDefusePerPlayer(p string) int {
	if playersdefuses == -1 {
		c := cmd.NewCommand(fmt.Sprintf(`ag -i "%s.+<STEAM.+ triggered \"Defused_The_Bomb\"" %s |wc -l`, p, testfile))

		err := c.Execute()
		if err != nil {
			panic(err.Error())
		}

		cs := strings.Trim(c.Stdout(), " \n")
		count, err := strconv.Atoi(cs)
		if err != nil {
			panic(err.Error())
		}
		playersdefuses = count
	}
	return playersdefuses
}

func TestDefusesPerPlayerByDB(t *testing.T) {
	var dbcount int
	player := "Jagger"
	filecount := countDefusePerPlayer(player)

	row := store.Db.QueryRow(context.Background(), "select count(scoreaction.id) from scoreaction "+
		"left join players on actor=players.id "+
		"where actiontype='defuse' and players.initialname=$1;", player)
	err := row.Scan(&dbcount)
	if err != nil {
		panic(err.Error())
	}

	if filecount != dbcount {
		t.Errorf("%s failed: filecount %d != dbcount %d", t.Name(), filecount, dbcount)
	}
}

var alldefuses int = -1

func countAllDefuses() int {
	if alldefuses == -1 {
		c := cmd.NewCommand(fmt.Sprintf(`ag -i ".+<STEAM.+ triggered \"Defused_The_Bomb\"" %s |wc -l`, testfile))

		err := c.Execute()
		if err != nil {
			panic(err.Error())
		}

		cs := strings.Trim(c.Stdout(), " \n")
		count, err := strconv.Atoi(cs)
		if err != nil {
			panic(err.Error())
		}
		alldefuses = count
	}
	return alldefuses
}

func TestAllDefusesByDB(t *testing.T) {
	var dbcount int
	filecount := countAllDefuses()

	row := store.Db.QueryRow(context.Background(), "select count(scoreaction.id) from scoreaction "+
		"left join players on actor=players.id "+
		"where actiontype='defuse';")

	err := row.Scan(&dbcount)
	if err != nil {
		panic(err.Error())
	}

	if filecount != dbcount {
		t.Errorf("%s failed: filecount %d != dbcount %d", t.Name(), filecount, dbcount)
	}
}
