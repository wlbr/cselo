package events

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/commander-cli/cmd"
)

var allbombings int = -1

func countAllBombings() int {
	if allbombings == -1 {
		c := cmd.NewCommand(fmt.Sprintf(`ag ".+ triggered \"SFUI_Notice_Target_Bombed\"" %s |wc -l`, testfile), cmd.WithInheritedEnvironment(nil))

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

func TestAllBombingsInMemory(t *testing.T) {
	filecount := countAllBombings()

	count := len(counter.Allbombings)
	if filecount != len(counter.Allbombings) {
		t.Errorf("%s failed: filecount %d != count %d", t.Name(), filecount, count)
	}
}
