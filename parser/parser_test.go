package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log-parser/match"
	"testing"
)

func Test_parseLog(t *testing.T) {
	type args struct {
		filepath string
	}
	tests := []struct {
		name        string
		args        args
		wantMatches []*match.Match
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "should return the expected matches information from the log file",
			args: args{
				filepath: "./testfiles/qgames_aborted_match.log",
			},
			wantMatches: []*match.Match{
				{
					TotalKills: 0,
					Players:    []string{"Isgalamido"},
					Kills: map[string]int{
						"Isgalamido": 0,
					},
					KillsByMeans: map[string]int{},
					PlayersInGame: map[string]bool{
						"Isgalamido": true,
					},
					Done:       true,
					InProgress: false,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "should return the expected matches information from the log file",
			args: args{
				filepath: "./testfiles/qgames_three_matches.log",
			},
			wantMatches: []*match.Match{
				{
					TotalKills: 0,
					Players:    []string{"Isgalamido"},
					Kills: map[string]int{
						"Isgalamido": 0,
					},
					KillsByMeans: map[string]int{},
					PlayersInGame: map[string]bool{
						"Isgalamido": true,
					},
					Done:       true,
					InProgress: false,
				},
				{
					TotalKills: 11,
					Players:    []string{"Isgalamido", "Dono da Bola", "Mocinha"},
					Kills: map[string]int{
						"Isgalamido":   -7,
						"Dono da Bola": 0,
						"Mocinha":      0,
					},
					KillsByMeans: map[string]int{
						"MOD_TRIGGER_HURT":  7,
						"MOD_ROCKET_SPLASH": 3,
						"MOD_FALLING":       1,
					},
					PlayersInGame: map[string]bool{
						"Isgalamido":   true,
						"Dono da Bola": true,
						"Mocinha":      true,
					},
					Done:       true,
					InProgress: false,
				},
				{
					TotalKills: 4,
					Players:    []string{"Dono da Bola", "Mocinha", "Isgalamido", "Zeh"},
					Kills: map[string]int{
						"Dono da Bola": -1,
						"Mocinha":      0,
						"Isgalamido":   1,
						"Zeh":          -2,
					},
					KillsByMeans: map[string]int{
						"MOD_ROCKET":       1,
						"MOD_TRIGGER_HURT": 2,
						"MOD_FALLING":      1,
					},
					PlayersInGame: map[string]bool{
						"Dono da Bola": true,
						"Mocinha":      true,
						"Isgalamido":   true,
						"Zeh":          true,
					},
					Done:       true,
					InProgress: false,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "should return the expected matches information from the log file",
			args: args{
				filepath: "./testfiles/qgames_complete_match.log",
			},
			wantMatches: []*match.Match{
				{
					TotalKills: 131,
					Players: []string{
						"Isgalamido",
						"Oootsimo",
						"Dono da Bola",
						"Assasinu Credi",
						"Zeh",
						"Mal",
					},
					Kills: map[string]int{
						"Isgalamido":     17,
						"Oootsimo":       21,
						"Dono da Bola":   12,
						"Assasinu Credi": 16,
						"Zeh":            19,
						"Mal":            6,
					},
					KillsByMeans: map[string]int{
						"MOD_ROCKET":        37,
						"MOD_TRIGGER_HURT":  14,
						"MOD_RAILGUN":       9,
						"MOD_ROCKET_SPLASH": 60,
						"MOD_MACHINEGUN":    4,
						"MOD_SHOTGUN":       4,
						"MOD_FALLING":       3,
					},
					PlayersInGame: map[string]bool{
						"Isgalamido":     true,
						"Oootsimo":       true,
						"Dono da Bola":   true,
						"Assasinu Credi": true,
						"Zeh":            true,
						"Mal":            true,
					},
					Done:       true,
					InProgress: false,
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLog(tt.args.filepath)
			if !tt.wantErr(t, err, fmt.Sprintf("parseLog(%v)", tt.args.filepath)) {
				return
			}

			for i, wantMatch := range tt.wantMatches {
				gotMatch := got[i]

				assert.Equal(t, wantMatch.TotalKills, gotMatch.TotalKills)
				assert.Equal(t, len(wantMatch.Players), len(gotMatch.Players))
				assert.Equal(t, wantMatch.InProgress, gotMatch.InProgress)
				assert.Equal(t, wantMatch.Done, gotMatch.Done)

				for k, v := range wantMatch.Kills {
					assert.Equal(t, v, gotMatch.Kills[k])
				}

				for k, v := range wantMatch.KillsByMeans {
					assert.Equal(t, v, gotMatch.KillsByMeans[k])
				}

				for k, v := range wantMatch.PlayersInGame {
					assert.Equal(t, v, gotMatch.PlayersInGame[k])
				}
			}
		})
	}
}
