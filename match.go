package main

type Match struct {
	TotalKills int      `json:"total_kills"`
	Players    []string `json:"players"`
	Kills      any      `json:"kills"`
	Done       bool
}

func NewMatch() *Match {
	return &Match{
		TotalKills: 0,
		Players:    make([]string, 0),
		Kills:      nil,
	}
}

func (m *Match) AddKillStats(player string) {

}
