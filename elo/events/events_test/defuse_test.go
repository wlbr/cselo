package events

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/commander-cli/cmd"
)

var playersdefuses int = -1

func countDefusePerPlayer(p string) int {
	if playersdefuses == -1 {
		c := cmd.NewCommand(fmt.Sprintf(`ag -i "%s.+<(\[U:|STEAM).+ triggered \"Defused_The_Bomb\"" %s |wc -l`, p, testfile), cmd.WithInheritedEnvironment(nil))

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

func TestDefusesPerPlayerInMemory(t *testing.T) {
	filecount := countDefusePerPlayer(player)

	count := len(counter.Playersdefuses)
	if filecount != count {
		t.Errorf("%s failed: filecount %d != count %d", t.Name(), filecount, count)
	}
}

var alldefuses int = -1

func countAllDefuses() int {
	if alldefuses == -1 {
		c := cmd.NewCommand(fmt.Sprintf(`ag -i ".+<(\[U:|STEAM).+ triggered \"Defused_The_Bomb\"" %s |wc -l`, testfile), cmd.WithInheritedEnvironment(nil))

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

func TestAllDefusesInMemory(t *testing.T) {
	filecount := countAllDefuses()

	count := len(counter.Alldefuses)
	if filecount != count {
		t.Errorf("%s failed: filecount %d != count %d", t.Name(), filecount, count)
	}
}
