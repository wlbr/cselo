package events

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/commander-cli/cmd"
)

func countKillsPerPlayer(p string) int {
	c := cmd.NewCommand(fmt.Sprintf(`ag -i "%s.+<STEAM.+killed.+<STEAM_" %s |wc -l`, p, testfile))

	err := c.Execute()
	if err != nil {
		panic(err.Error())
	}

	cs := strings.Trim(c.Stdout(), " \n")
	killcount, err := strconv.Atoi(cs)
	if err != nil {
		panic(err.Error())
	}
	return killcount
}

func TestKillsPerPlayer(t *testing.T) {
	var dbkillcount int
	player := "Jagger"
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

func countAllKills() int {
	c := cmd.NewCommand(fmt.Sprintf(`ag -i ".+<STEAM.+killed.+<STEAM_" %s |wc -l`, testfile))

	err := c.Execute()
	if err != nil {
		panic(err.Error())
	}

	cs := strings.Trim(c.Stdout(), " \n")
	killcount, err := strconv.Atoi(cs)
	if err != nil {
		panic(err.Error())
	}
	return killcount
}

func TestAllKills(t *testing.T) {
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
