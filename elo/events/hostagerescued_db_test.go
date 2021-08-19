package events

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/commander-cli/cmd"
)

var playersrescues int = -1

func countHostageRescuePerPlayer(p string) int {
	if playersrescues == -1 {
		c := cmd.NewCommand(fmt.Sprintf(`ag -i "%s.+<STEAM.+ triggered \"Rescued_A_Hostage\"" %s |wc -l`, p, testfile))

		err := c.Execute()
		if err != nil {
			panic(err.Error())
		}

		cs := strings.Trim(c.Stdout(), " \n")
		count, err := strconv.Atoi(cs)
		if err != nil {
			panic(err.Error())
		}
		playersrescues = count
	}
	return playersrescues
}

func TestHostageRescuesPerPlayerByDB(t *testing.T) {
	var dbcount int
	player := "Jagger"
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

var allrescues int = -1

func countAllHostageRescues() int {
	if allrescues == -1 {
		c := cmd.NewCommand(fmt.Sprintf(`ag -i ".+<STEAM.+ triggered \"Rescued_A_Hostage\"" %s |wc -l`, testfile))

		err := c.Execute()
		if err != nil {
			panic(err.Error())
		}

		cs := strings.Trim(c.Stdout(), " \n")
		count, err := strconv.Atoi(cs)
		if err != nil {
			panic(err.Error())
		}
		allrescues = count
	}
	return allrescues
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
