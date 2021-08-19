package events

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/commander-cli/cmd"
)

var allbombings int = -1

func countAllBombings() int {
	if allbombings == -1 {
		c := cmd.NewCommand(fmt.Sprintf(`ag ".+ triggered \"SFUI_Notice_Target_Bombed\"" %s |wc -l`, testfile))

		err := c.Execute()
		if err != nil {
			panic(err.Error())
		}

		cs := strings.Trim(c.Stdout(), " \n")
		count, err := strconv.Atoi(cs)
		if err != nil {
			panic(err.Error())
		}
		allassists = count
	}
	return allassists
}

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
