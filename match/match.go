package match

const (
	world = "<world>"
)

type Match struct {
	TotalKills    int            `json:"total_kills"`
	Players       []string       `json:"players"`
	Kills         map[string]int `json:"kills"`
	KillsByMeans  map[string]int `json:"kills_by_means"`
	PlayersInGame map[string]bool
	Done          bool
	InProgress    bool
}

func NewMatch() *Match {
	return &Match{
		TotalKills:    0,
		Players:       make([]string, 0),
		Kills:         map[string]int{},
		KillsByMeans:  map[string]int{},
		PlayersInGame: map[string]bool{},
		Done:          false,
		InProgress:    false,
	}
}

func (m *Match) AddKillStats(player string) {
	_, ok := m.Kills[player]
	if !ok {
		m.Kills[player] = 0
	}
}

func (m *Match) AddKillAndMeans(killer, killed, reason string) {
	m.KillsByMeans[reason]++
	m.TotalKills++

	if killer == killed {
		return
	}

	if killer == world {
		if m.Kills[killed] > 0 {
			m.Kills[killed]--
		}
	} else {
		m.Kills[killer]++
	}
}
