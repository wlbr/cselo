package elo

import "regexp"

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
		return true
	}
	return false
}
