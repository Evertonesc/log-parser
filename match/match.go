package match

const (
	world = "<world>"
)

type (
	Match struct {
		TotalKills    int             `json:"total_kills"`
		Players       []string        `json:"players"`
		Kills         map[string]int  `json:"kills"`
		KillsByMeans  map[string]int  `json:"-"`
		PlayersInGame map[string]bool `json:"-"`
		Done          bool            `json:"-"`
		InProgress    bool            `json:"-"`
	}

	Summary struct {
		KillsByMeans map[string]int `json:"kills_by_means"`
	}
)

func NewMatch() *Match {
	return &Match{
		TotalKills:    0,
		Players:       make([]string, 0),
		Kills:         make(map[string]int),
		KillsByMeans:  make(map[string]int),
		PlayersInGame: make(map[string]bool),
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
		m.Kills[killed]--
	} else {
		m.Kills[killer]++
	}
}
