package events

import (
	"os"
	"testing"
	"time"

	"github.com/wlbr/cs-elo/elo"
	"github.com/wlbr/cs-elo/elo/sources/postgresql"
)

var config *elo.Config
var store *postgresql.Postgres
var testfile string

func TestMain(m *testing.M) {
	testsetup()
	code := m.Run()
	testteardown()
	os.Exit(code)
}

func testsetup() {
	testfile = "../../data/test.log"

	config = new(elo.Config)
	config.Initialize("Test", time.Now().Format(time.ANSIC))

	// db tests
	var err error
	store, err = postgresql.NewPostgres(config)
	if err != nil {
		panic(err.Error())
	}

}

func testteardown() {
	defer config.CleanUp()

}
