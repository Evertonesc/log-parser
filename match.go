package main

type Match struct {
	TotalKills   int            `json:"total_kills"`
	Players      []string       `json:"players"`
	Kills        map[string]int `json:"kills"`
	KillsByMeans map[string]int `json:"kills_by_means"`
	Done         bool
	InProgress   bool
}

func NewMatch() *Match {
	return &Match{
		TotalKills:   0,
		Players:      make([]string, 0),
		Kills:        map[string]int{},
		KillsByMeans: map[string]int{},
		Done:         false,
	}
}

func (m *Match) AddKillStats(player string) {
	m.Kills = map[string]int{
		player: 0,
	}
}
