package events

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/commander-cli/cmd"
)

var playerskills int = -1

func countKillsPerPlayer(p string) int {
	if playerskills == -1 {
		c := cmd.NewCommand(fmt.Sprintf(`ag -i "%s.+<STEAM.+killed.+<STEAM_" %s |wc -l`, p, testfile))

		err := c.Execute()
		if err != nil {
			panic(err.Error())
		}

		cs := strings.Trim(c.Stdout(), " \n")
		count, err := strconv.Atoi(cs)
		if err != nil {
			panic(err.Error())
		}
		playerskills = count
	}
	return playerskills
}

func TestKillsPerPlayerInMemory(t *testing.T) {
	killcount := countKillsPerPlayer(player)

	count := len(counter.PlayersKills)
	if killcount != count {
		t.Errorf("%s failed: filecount %d != count %d", t.Name(), killcount, count)
	}
}

var allkills int = -1

func countAllKills() int {
	if allkills == -1 {
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
		allkills = killcount
	}
	return allkills
}

func TestAllKillsInMemory(t *testing.T) {
	killcount := countAllKills()

	count := len(counter.AllKills)
	if killcount != count {
		t.Errorf("%s failed: filecount %d != count %d", t.Name(), killcount, count)
	}
}
