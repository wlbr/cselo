package events

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/commander-cli/cmd"
)

var matchends int = -1

func countMatchEnds() int {
	if matchends == -1 {
		c := cmd.NewCommand(fmt.Sprintf(`ag -i "Game Over:" %s |wc -l`, testfile))

		err := c.Execute()
		if err != nil {
			panic(err.Error())
		}

		cs := strings.Trim(c.Stdout(), " \n")
		count, err := strconv.Atoi(cs)
		if err != nil {
			panic(err.Error())
		}
		matchends = count
	}
	return matchends
}

func TestMatchEndsInMemory(t *testing.T) {
	matchendscount := countMatchEnds()

	count := len(counter.Matches)
	if matchendscount != count {
		t.Errorf("%s failed: filecount %d != count %d", t.Name(), matchendscount, count)
	}
}
