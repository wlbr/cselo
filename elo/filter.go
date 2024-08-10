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

func CheckFilter(em Emitter, m string) bool {
	for _, f := range em.GetFilters() {
		if f.Test(m) {
			return true
		}
	}
	return false
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

var steamidpendingrex = regexp.MustCompile(`<STEAM_ID_PENDING>`)

type SteamIdPendingFilter struct {
}

func (f *SteamIdPendingFilter) String() string {
	return "Filters all events with STEAM_ID_PENDING players involved."
}

func (f *SteamIdPendingFilter) Test(message string) bool {
	if steamidpendingrex.MatchString(message) {
		log.Info("SteamIdPendingFilter: Filtered message '%s'", message)
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
