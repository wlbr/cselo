package events

import (
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

func TestHostageRescuesPerPlayerInMemory(t *testing.T) {
	filecount := countHostageRescuePerPlayer(player)

	count := len(counter.Playersrescues)

	if filecount != count {
		t.Errorf("%s failed: filecount %d != count %d", t.Name(), filecount, count)
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

func TestAllHostageRescuesInMemory(t *testing.T) {
	filecount := countAllHostageRescues()

	count := len(counter.Allrescues)

	if filecount != count {
		t.Errorf("%s failed: filecount %d != dbcount %d", t.Name(), filecount, count)
	}
}
