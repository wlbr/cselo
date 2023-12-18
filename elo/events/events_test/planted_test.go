package events

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/commander-cli/cmd"
)

func countPlantingsPerPlayer(p string) int {
	c := cmd.NewCommand(fmt.Sprintf(`ag -i "%s.+<(\[U:|STEAM).+ triggered \"Planted_The_Bomb\"" %s |wc -l`, p, testfile), cmd.WithInheritedEnvironment(nil))

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

func TestPlantingsPerPlayerInMemory(t *testing.T) {
	filecount := countPlantingsPerPlayer(player)

	count := len(counter.Playersplantings)
	if filecount != count {
		t.Errorf("%s failed: filecount %d != count %d", t.Name(), filecount, count)
	}
}

func countAllPlantings() int {
	c := cmd.NewCommand(fmt.Sprintf(`ag -i ".+<(\[U:|STEAM).+ triggered \"Planted_The_Bomb\"" %s |wc -l`, testfile), cmd.WithInheritedEnvironment(nil))

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

func TestAllPlantingsInMemory(t *testing.T) {
	filecount := countAllPlantings()

	count := len(counter.Allplantings)
	if filecount != count {
		t.Errorf("%s failed: filecount %d != count %d", t.Name(), filecount, count)
	}
}
