package events

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/commander-cli/cmd"
)

func countAssistsPerPlayer(p string) int {
	c := cmd.NewCommand(fmt.Sprintf(`ag -i "%s.+<STEAM.+assisted.+<STEAM_" %s |wc -l`, p, testfile))

	err := c.Execute()
	if err != nil {
		panic(err.Error())
	}

	cs := strings.Trim(c.Stdout(), " \n")
	count, err := strconv.Atoi(cs)
	if err != nil {
		panic(err.Error())
	}
	return count
}

func TestAssitsPerPlayer(t *testing.T) {
	var dbcount int
	player := "Jagger"
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

func countAllAssists() int {
	c := cmd.NewCommand(fmt.Sprintf(`ag -i ".+<STEAM.+assisted.+<STEAM_" %s |wc -l`, testfile))

	err := c.Execute()
	if err != nil {
		panic(err.Error())
	}

	cs := strings.Trim(c.Stdout(), " \n")
	count, err := strconv.Atoi(cs)
	if err != nil {
		panic(err.Error())
	}
	return count
}

func TestAllAssits(t *testing.T) {
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
