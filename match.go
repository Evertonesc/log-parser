package main

type Match struct {
	TotalKills int            `json:"total_kills"`
	Players    []string       `json:"players"`
	Kills      map[string]int `json:"kills"`
	Done       bool
}

func NewMatch() *Match {
	return &Match{}
}

func (m *Match) InitPlayers() {
	m.Players = make([]string, 0)
}

func (m *Match) AddKillStats(player string) {
	if m.Kills == nil {
		m.Kills = make(map[string]int)
	}

	m.Kills = map[string]int{
		player: 0,
	}
}
