package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	match "log-parser/match"
	"testing"
)

func TestInitGameHandler_Handle(t *testing.T) {
	type args struct {
		logLine string
		match   *match.Match
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should successfully parse a init game log entry",
			args: args{
				logLine: "0:00 InitGame: \\\\\\\\sv_floodProtect\\\\\\\\1\\\\\\\\sv_maxPing\\\\\\\\0\\\\\\\\sv_minPing\\\\\\\\0\\\\\\\\sv_maxRate\\\\\\\\10000\\\\\\\\sv_minRate\\\\\\\\0\\\\\\\\sv_hostname\\\\\\\\Code Miner Server\\\\\\\\g_gametype\\\\\\\\0\\\\\\\\sv_privateClients\\\\\\\\2\\\\\\\\sv_maxclients\\\\\\\\16\\\\\\\\sv_allowDownload\\\\\\\\0\\\\\\\\dmflags\\\\\\\\0\\\\\\\\fraglimit\\\\\\\\20\\\\\\\\timelimit\\\\\\\\15\\\\\\\\g_maxGameClients\\\\\\\\0\\\\\\\\capturelimit\\\\\\\\8\\\\\\\\version\\\\\\\\ioq3 1.36 linux-x86_64 Apr 12 2009\\\\\\\\protocol\\\\\\\\68\\\\\\\\mapname\\\\\\\\q3dm17\\\\\\\\gamename\\\\\\\\baseq3\\\\\\\\g_needpass\\\\\\\\0\\n",
				match:   match.NewMatch(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewInitGameHandler()
			tt.wantErr(t, h.Handle(tt.args.logLine, tt.args.match), fmt.Sprintf("Handle(%v, %v)", tt.args.logLine, tt.args.match))

			m := tt.args.match

			assert.Equal(t, true, m.InProgress)
		})
	}
}

func TestAddPlayerHandler_Handle(t *testing.T) {
	type args struct {
		logLine string
		match   func() *match.Match
	}
	tests := []struct {
		name      string
		args      args
		wantMatch *match.Match
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "should successfully parse a client user info changed log entry and update the match data with a new player",
			args: args{
				logLine: " 20:34 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0\\model\\xian/default\\hmodel\\xian/default\\g_redteam\\\\g_blueteam\\\\c1\\4\\c2\\5\\hc\\100\\w\\0\\l\\0\\tt\\0\\tl\\0",
				match: func() *match.Match {
					m := match.NewMatch()
					m.InProgress = true

					return m
				},
			},
			wantMatch: &match.Match{
				TotalKills: 0,
				Players:    []string{"Isgalamido"},
				Kills: map[string]int{
					"Isgalamido": 0,
				},
				KillsByMeans: nil,
				PlayersInGame: map[string]bool{
					"Isgalamido": true,
				},
				Done:       false,
				InProgress: true,
			},
			wantErr: assert.NoError,
		},
		{
			name: "should successfully parse a client user info changed log entry and update the match data with a new player",
			args: args{
				logLine: " 20:34 ClientUserinfoChanged: 2 n\\Mocinha\\t\\0\\model\\xian/default\\hmodel\\xian/default\\g_redteam\\\\g_blueteam\\\\c1\\4\\c2\\5\\hc\\100\\w\\0\\l\\0\\tt\\0\\tl\\0",
				match: func() *match.Match {
					return &match.Match{
						TotalKills: 0,
						Players:    []string{"Isgalamido"},
						Kills: map[string]int{
							"Isgalamido": 0,
						},
						KillsByMeans: nil,
						PlayersInGame: map[string]bool{
							"Isgalamido": true,
						},
						Done:       false,
						InProgress: true,
					}
				},
			},
			wantMatch: &match.Match{
				TotalKills: 0,
				Players:    []string{"Isgalamido", "Mocinha"},
				Kills: map[string]int{
					"Isgalamido": 0,
					"Mocinha":    0,
				},
				KillsByMeans: nil,
				PlayersInGame: map[string]bool{
					"Isgalamido": true,
					"Mocinha":    true,
				},
				Done:       false,
				InProgress: true,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewAddPlayerHandler()
			m := tt.args.match()
			tt.wantErr(t, h.Handle(tt.args.logLine, m), fmt.Sprintf("Handle(%v, %v)", tt.args.logLine, m))

			assert.Equal(t, len(tt.wantMatch.Players), len(m.Players))
			assert.Equal(t, tt.wantMatch.Players[0], m.Players[0])

			for k, v := range tt.wantMatch.Kills {
				assert.Equal(t, v, m.Kills[k])
			}

			for k, v := range tt.wantMatch.PlayersInGame {
				assert.Equal(t, v, m.PlayersInGame[k])
			}
		})
	}
}

func TestKillDetailsHandler_Handle(t *testing.T) {
	type args struct {
		logLine string
		match   func() *match.Match
	}
	tests := []struct {
		name      string
		args      args
		wantMatch *match.Match
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "should successfully parse the kill log entry and add its information to the match",
			args: args{
				logLine: " 20:54 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT",
				match: func() *match.Match {
					return &match.Match{
						TotalKills: 5,
						Players:    []string{"Isgalamido", "Bruce Wayne"},
						Kills: map[string]int{
							"Isgalamido":  3,
							"Bruce Wayne": 2,
						},
						KillsByMeans: map[string]int{
							"MOD_TRIGGER_HURT":  1,
							"MOD_ROCKET_SPLASH": 2,
							"MOD_FALLING":       2,
						},
						PlayersInGame: map[string]bool{
							"Isgalamido":  true,
							"Bruce Wayne": true,
						},
						Done:       false,
						InProgress: true,
					}

				},
			},
			wantMatch: &match.Match{
				TotalKills: 6,
				Players:    []string{"Isgalamido", "Bruce Wayne"},
				Kills: map[string]int{
					"Isgalamido":  2,
					"Bruce Wayne": 2,
				},
				KillsByMeans: map[string]int{
					"MOD_TRIGGER_HURT":  2,
					"MOD_ROCKET_SPLASH": 2,
					"MOD_FALLING":       2,
				},
				PlayersInGame: map[string]bool{
					"Isgalamido":  true,
					"Bruce Wayne": true,
				},
				Done:       false,
				InProgress: true,
			},
			wantErr: assert.NoError,
		},
		{
			name: "should successfully parse the kill log entry and add its information to the match",
			args: args{
				logLine: "  1:41 Kill: 1022 2 19: <world> killed Dono da Bola by MOD_FALLING",
				match: func() *match.Match {
					return &match.Match{
						TotalKills: 5,
						Players:    []string{"Isgalamido", "Bruce Wayne"},
						Kills: map[string]int{
							"Isgalamido":   3,
							"Dono da Bola": 2,
						},
						KillsByMeans: map[string]int{
							"MOD_TRIGGER_HURT":  1,
							"MOD_ROCKET_SPLASH": 2,
							"MOD_FALLING":       2,
						},
						PlayersInGame: map[string]bool{
							"Isgalamido":   true,
							"Dono da Bola": true,
						},
						Done:       false,
						InProgress: true,
					}

				},
			},
			wantMatch: &match.Match{
				TotalKills: 6,
				Players:    []string{"Isgalamido", "Bruce Wayne"},
				Kills: map[string]int{
					"Isgalamido":   3,
					"Dono da Bola": 1,
				},
				KillsByMeans: map[string]int{
					"MOD_TRIGGER_HURT":  1,
					"MOD_ROCKET_SPLASH": 2,
					"MOD_FALLING":       3,
				},
				PlayersInGame: map[string]bool{
					"Isgalamido":   true,
					"Dono da Bola": true,
				},
				Done:       false,
				InProgress: true,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewKillDetailsHandler()
			m := tt.args.match()
			tt.wantErr(t, h.Handle(tt.args.logLine, m), fmt.Sprintf("Handle(%v, %v)", tt.args.logLine, m))

			assert.Equal(t, len(tt.wantMatch.Players), len(m.Players))
			assert.Equal(t, tt.wantMatch.Players[0], m.Players[0])

			for k, v := range tt.wantMatch.Kills {
				assert.Equal(t, v, m.Kills[k])
			}

			for k, v := range tt.wantMatch.PlayersInGame {
				assert.Equal(t, v, m.PlayersInGame[k])
			}

			for k, v := range tt.wantMatch.KillsByMeans {
				assert.Equal(t, v, m.KillsByMeans[k])
			}
		})
	}
}

func TestEndGameHandler_Handle(t *testing.T) {
	type args struct {
		logLine string
		match   *match.Match
	}
	tests := []struct {
		name      string
		args      args
		wantMatch *match.Match
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "should successfully parse a shutdown game log entry and set the done match stats to true and in progress to false",
			args: args{
				logLine: " 20:37 ShutdownGame:",
				match: &match.Match{
					TotalKills: 5,
					Players:    []string{"Isgalamido", "Bruce Wayne"},
					Kills: map[string]int{
						"Isgalamido":  3,
						"Bruce Wayne": 2,
					},
					KillsByMeans: map[string]int{
						"MOD_TRIGGER_HURT":  1,
						"MOD_ROCKET_SPLASH": 2,
						"MOD_FALLING":       2,
					},
					PlayersInGame: map[string]bool{
						"Isgalamido":  true,
						"Bruce Wayne": true,
					},
					Done:       false,
					InProgress: true,
				},
			},
			wantMatch: &match.Match{
				TotalKills: 5,
				Players:    []string{"Isgalamido", "Bruce Wayne"},
				Kills: map[string]int{
					"Isgalamido":  3,
					"Bruce Wayne": 2,
				},
				KillsByMeans: map[string]int{
					"MOD_TRIGGER_HURT":  1,
					"MOD_ROCKET_SPLASH": 2,
					"MOD_FALLING":       2,
				},
				PlayersInGame: map[string]bool{
					"Isgalamido":  true,
					"Bruce Wayne": true,
				},
				Done:       true,
				InProgress: false,
			},
			wantErr: assert.NoError,
		},
		{
			name: "should successfully parse a log entry that finishes the game by an unknown reason and finish the match",
			args: args{
				logLine: "26  0:00 ------------------------------------------------------------",
				match: &match.Match{
					TotalKills: 5,
					Players:    []string{"Isgalamido", "Bruce Wayne"},
					Kills: map[string]int{
						"Isgalamido":  3,
						"Bruce Wayne": 2,
					},
					KillsByMeans: map[string]int{
						"MOD_TRIGGER_HURT":  1,
						"MOD_ROCKET_SPLASH": 2,
						"MOD_FALLING":       2,
					},
					PlayersInGame: map[string]bool{
						"Isgalamido":  true,
						"Bruce Wayne": true,
					},
					Done:       false,
					InProgress: true,
				},
			},
			wantMatch: &match.Match{
				TotalKills: 5,
				Players:    []string{"Isgalamido", "Bruce Wayne"},
				Kills: map[string]int{
					"Isgalamido":  3,
					"Bruce Wayne": 2,
				},
				KillsByMeans: map[string]int{
					"MOD_TRIGGER_HURT":  1,
					"MOD_ROCKET_SPLASH": 2,
					"MOD_FALLING":       2,
				},
				PlayersInGame: map[string]bool{
					"Isgalamido":  true,
					"Bruce Wayne": true,
				},
				Done:       true,
				InProgress: false,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewEndGameHandler()
			tt.wantErr(t, h.Handle(tt.args.logLine, tt.args.match), fmt.Sprintf("Handle(%v, %v)", tt.args.logLine, tt.args.match))

			assert.Equal(t, len(tt.wantMatch.Players), len(tt.args.match.Players))

			wantedPlayersCount := len(tt.wantMatch.Players)
			for i := 0; i < wantedPlayersCount; i++ {
				assert.Equal(t, tt.wantMatch.Players[i], tt.args.match.Players[i])
			}

			for k, v := range tt.wantMatch.Kills {
				assert.Equal(t, v, tt.args.match.Kills[k])
			}

			for k, v := range tt.wantMatch.PlayersInGame {
				assert.Equal(t, v, tt.args.match.PlayersInGame[k])
			}

			for k, v := range tt.wantMatch.KillsByMeans {
				assert.Equal(t, v, tt.args.match.KillsByMeans[k])
			}
		})
	}
}
