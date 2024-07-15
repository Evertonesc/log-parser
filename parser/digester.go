package parser

import (
	"log-parser/match"
	"regexp"
)

var (
	initGameRe             = regexp.MustCompile(`^\s*\d{1,3}:\d{2}\s+InitGame:.*$`)
	clientUserInfoRe       = regexp.MustCompile(`ClientUserinfoChanged:\s+\d+\s+n\\([^\\]+)`)
	killDetailsRe          = regexp.MustCompile(`\s*\d{1,2}:\d{2}\s+Kill: \d+ \d+ \d+: ([^ ]+) killed ([^ ]+(?: [^ ]+)*) by ([^ ]+)`)
	shutDownGameRe         = regexp.MustCompile(`^\s*\d{1,3}:\d{2}\s+ShutdownGame:$`)
	unknownReasonEndGameRe = regexp.MustCompile(`^.*\d+\s+0:00`)
)

type LogDigesterHandler interface {
	Handle(logLine string, match *match.Match) error
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

	KillDetailsHandler struct {
		generalLogDigesterHandler
	}

	EndGameHandler struct {
		generalLogDigesterHandler
	}
)

func (h *generalLogDigesterHandler) SetNext(handler LogDigesterHandler) {
	h.Next = handler
}

func (h *generalLogDigesterHandler) handleNext(logLine string, match *match.Match) error {
	if h.Next != nil {
		return h.Next.Handle(logLine, match)
	}

	return nil
}

func NewInitGameHandler() *InitGameHandler {
	return &InitGameHandler{}
}

func (h *InitGameHandler) Handle(logLine string, match *match.Match) error {
	if initGameRe.MatchString(logLine) {
		if !match.InProgress {
			match.InProgress = true

			return nil
		}

		match.Done = true

		return nil
	}

	return h.handleNext(logLine, match)
}

func NewAddPlayerHandler() *AddPlayerHandler {
	return &AddPlayerHandler{}
}

func (h *AddPlayerHandler) Handle(logLine string, match *match.Match) error {
	values := clientUserInfoRe.FindStringSubmatch(logLine)
	if len(values) > 0 {
		player := values[1]

		_, ok := match.PlayersInGame[player]
		if !ok {
			match.Players = append(match.Players, player)
			match.PlayersInGame[player] = true
			match.AddKillStats(player)
		}
	}

	return h.handleNext(logLine, match)
}

func NewKillDetailsHandler() *KillDetailsHandler {
	return &KillDetailsHandler{}
}

func (h *KillDetailsHandler) Handle(logLine string, match *match.Match) error {
	killerPiece, killedPlayerPiece, reasonPiece := 1, 2, 3
	matches := killDetailsRe.FindStringSubmatch(logLine)
	if len(matches) > 3 {
		match.AddKillAndMeans(
			matches[killerPiece],
			matches[killedPlayerPiece],
			matches[reasonPiece],
		)

		return nil
	}

	return h.handleNext(logLine, match)
}

func NewEndGameHandler() *EndGameHandler {
	return &EndGameHandler{}
}

func (h *EndGameHandler) Handle(logLine string, match *match.Match) error {
	if shutDownGameRe.MatchString(logLine) || unknownReasonEndGameRe.MatchString(logLine) {
		match.InProgress = false
		match.Done = true
	}

	return nil
}

func LoadLogsDigester() LogDigesterHandler {
	endGameHandler := NewEndGameHandler()

	killDetailsHandler := NewKillDetailsHandler()
	killDetailsHandler.SetNext(endGameHandler)

	addPlayerHandler := NewAddPlayerHandler()
	addPlayerHandler.SetNext(killDetailsHandler)

	initGameHandler := NewInitGameHandler()
	initGameHandler.SetNext(addPlayerHandler)

	return initGameHandler
}
