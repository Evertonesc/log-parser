package match

import "testing"

func TestMatch_AddKillAndMeans(t *testing.T) {
	type fields struct {
		TotalKills   int
		Players      []string
		Kills        map[string]int
		KillsByMeans map[string]int
		Done         bool
	}
	type args struct {
		killer string
		killed string
		reason string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantMatch *Match
	}{
		{
			name: "should insert the kill information to the match",
			fields: fields{
				TotalKills: 0,
				Players:    []string{"Isgalamido", "Bruce Wayne"},
				Kills: map[string]int{
					"Isgalamido":  0,
					"Bruce Wayne": 0,
				},
				KillsByMeans: map[string]int{},
			},
			args: args{
				killer: "Bruce Wayne",
				killed: "Isgalamido",
				reason: "MOD_ROCKET_SPLASH",
			},
			wantMatch: &Match{
				TotalKills: 1,
				Players:    []string{"Isgalamido", "Bruce Wayne"},
				Kills: map[string]int{
					"Isgalamido":  0,
					"Bruce Wayne": 1,
				},
				KillsByMeans: map[string]int{
					"MOD_ROCKET_SPLASH": 1,
				},
			},
		},
		{
			name: "should insert the kill information and decrease the killed player by world",
			fields: fields{
				TotalKills: 10,
				Players:    []string{"Isgalamido", "Bruce Wayne"},
				Kills: map[string]int{
					"Isgalamido":  8,
					"Bruce Wayne": 2,
				},
				KillsByMeans: map[string]int{
					"MOD_ROCKET_SPLASH": 1,
				},
			},
			args: args{
				killer: "<world>",
				killed: "Isgalamido",
				reason: "MOD_TRIGGER_HURT",
			},
			wantMatch: &Match{
				TotalKills: 11,
				Players:    []string{"Isgalamido", "Bruce Wayne"},
				Kills: map[string]int{
					"Isgalamido":  7,
					"Bruce Wayne": 2,
				},
				KillsByMeans: map[string]int{
					"MOD_ROCKET_SPLASH": 1,
					"MOD_TRIGGER_HURT":  1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Match{
				TotalKills:   tt.fields.TotalKills,
				Players:      tt.fields.Players,
				Kills:        tt.fields.Kills,
				KillsByMeans: tt.fields.KillsByMeans,
				Done:         tt.fields.Done,
			}
			m.AddKillAndMeans(tt.args.killer, tt.args.killed, tt.args.reason)
		})
	}
}
