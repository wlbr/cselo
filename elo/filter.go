package elo

import (
	"regexp"

	"github.com/wlbr/commons/log"
)

type Filter interface {
	String() string
	Test(string) bool
}

//================================

var allbotsrex = regexp.MustCompile(`<BOT>`)

type AllBotsFilter struct {
}

func (f *AllBotsFilter) String() string {
	return "Filters all events with bots involved."
}

func (f *AllBotsFilter) Test(message string) bool {
	if allbotsrex.MatchString(message) {
		log.Info("AllBotsFilter: Filtered message '%s'", message)
		return true
	}
	return false
}
