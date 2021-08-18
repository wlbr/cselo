package events

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/commander-cli/cmd"
)

func countHostageRescuePerPlayer(p string) int {
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
	return count
}

func TestHostageRescues(t *testing.T) {
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
