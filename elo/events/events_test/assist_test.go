package events

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/commander-cli/cmd"
)

var playersassists int = -1

func countAssistsPerPlayer(p string) int {
	if playersassists == -1 {
		c := cmd.NewCommand(fmt.Sprintf(`ag -i "%s.+<(\[U:|STEAM).+ assisted.+<(\[U:|STEAM)" %s |wc -l`, p, testfile), cmd.WithInheritedEnvironment(nil))

		err := c.Execute()
		if err != nil {
			panic(err.Error())
		}

		cs := strings.Trim(c.Stdout(), " \n")
		count, err := strconv.Atoi(cs)
		if err != nil {
			panic(err.Error())
		}
		playersassists = count
	}
	return playersassists
}

func TestAssistsPerPlayerInMemory(t *testing.T) {
	filecount := countAssistsPerPlayer(player)

	count := len(counter.Playersassists)
	if filecount != count {
		t.Errorf("%s failed: filecount %d != count %d", t.Name(), filecount, count)
	}
}

var allassists int = -1

func countAllAssists() int {
	if allassists == -1 {
		c := cmd.NewCommand(fmt.Sprintf(`ag -i ".+<(\[U:|STEAM).+ assisted.+<(\[U:|STEAM)" %s |wc -l`, testfile), cmd.WithInheritedEnvironment(nil))

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

func TestAllAssistsInMemory(t *testing.T) {
	filecount := countAllAssists()

	count := len(counter.Allassists)

	if filecount != count {
		t.Errorf("%s failed: filecount %d != count %d", t.Name(), filecount, count)
	}
}
