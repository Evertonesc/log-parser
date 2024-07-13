package main

import (
	"regexp"
)

var (
	initGameRe       = regexp.MustCompile(`^\s*\d+:\d+\s+InitGame:`)
	clientUserInfoRe = regexp.MustCompile(`ClientUserinfoChanged:\s+\d+\s+n\\([^\\]+)`)
)

type LogDigesterHandler interface {
	Handle(logLine string, match *Match) error
}

type (
	generalLogDigesterHandler struct {
		Next LogDigesterHandler
	}

	InitGameHandler struct {
		generalLogDigesterHandler
	}

	AddPlayerHandler struct {
		generalLogDigesterHandler
	}
)

func (h *generalLogDigesterHandler) SetNext(handler LogDigesterHandler) {
	h.Next = handler
}

func (h *generalLogDigesterHandler) handleNext(logLine string, match *Match) error {
	if h.Next != nil {
		return h.Next.Handle(logLine, match)
	}

	return nil
}

func NewInitGameHandler() *InitGameHandler {
	return &InitGameHandler{}
}

func (h *InitGameHandler) Handle(logLine string, match *Match) error {
	if initGameRe.MatchString(logLine) {
		match = NewMatch()
		return nil
	}

	return h.handleNext(logLine, match)
}

func NewAddPlayerHandler() *AddPlayerHandler {
	return &AddPlayerHandler{}
}

func (h *AddPlayerHandler) Handle(logLine string, match *Match) error {
	values := clientUserInfoRe.FindStringSubmatch(logLine)
	if len(values) > 0 {
		if len(match.Players) == 0 {
			match.Players = append(match.Players, values[1])
		} else {
			for _, player := range match.Players {
				if player != values[1] {
					match.Players = append(match.Players, values[1])
				}
			}
		}

	}

	return h.handleNext(logLine, match)
}

func LoadLogsDigester() LogDigesterHandler {
	addPlayerHandler := NewAddPlayerHandler()

	initGameHandler := NewInitGameHandler()
	initGameHandler.SetNext(addPlayerHandler)

	return initGameHandler
}
