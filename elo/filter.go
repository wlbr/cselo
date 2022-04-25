package elo

import (
	"regexp"
	"strings"

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
	if allbotsrex.MatchString(message) && !strings.Contains(message, "Punting bot, server is hibernating") {
		log.Info("AllBotsFilter: Filtered message '%s'", message)
		return true
	}
	return false
}

//================================

var unkownrex = regexp.MustCompile(`<unknown>`)

type UnknownFilter struct {
}

func (f *UnknownFilter) String() string {
	return "Filters all events with unknown players involved."
}

func (f *UnknownFilter) Test(message string) bool {
	if unkownrex.MatchString(message) {
		log.Info("UnknownFilter: Filtered message '%s'", message)
		return true
	}
	return false
}
