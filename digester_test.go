package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogsDigester(t *testing.T) {
	type args struct {
		logLine string
		match   *Match
	}
	tests := []struct {
		name      string
		args      args
		wantErr   error
		wantMatch *Match
	}{
		{
			name: "should parse the init game log entry",
			args: args{
				logLine: "0:00 InitGame: \\\\sv_floodProtect\\\\1\\\\sv_maxPing\\\\0\\\\sv_minPing\\\\0\\\\sv_maxRate\\\\10000\\\\sv_minRate\\\\0\\\\sv_hostname\\\\Code Miner Server\\\\g_gametype\\\\0\\\\sv_privateClients\\\\2\\\\sv_maxclients\\\\16\\\\sv_allowDownload\\\\0\\\\dmflags\\\\0\\\\fraglimit\\\\20\\\\timelimit\\\\15\\\\g_maxGameClients\\\\0\\\\capturelimit\\\\8\\\\version\\\\ioq3 1.36 linux-x86_64 Apr 12 2009\\\\protocol\\\\68\\\\mapname\\\\q3dm17\\\\gamename\\\\baseq3\\\\g_needpass\\\\0\n",
				match:   &Match{},
			},
			wantErr: nil,
			wantMatch: &Match{
				Players: []string{},
			},
		},
		{
			name: "should parse the client user info changed log entry adding a new player to the match",
			args: args{
				logLine: " 20:34 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0\\model\\xian/default\\hmodel\\xian/default\\g_redteam\\\\g_blueteam\\\\c1\\4\\c2\\5\\hc\\100\\w\\0\\l\\0\\tt\\0\\tl\\0",
				match:   &Match{},
			},
			wantMatch: &Match{
				Players: []string{"Isgalamido"},
				Kills: map[string]int{
					"Isgalamido": 0,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			logsDigester := LoadLogsDigester()
			err := logsDigester.Handle(tt.args.logLine, tt.args.match)

			if tt.wantErr != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("error mismtach, want %s, got %s", tt.wantErr.Error(), err.Error())
				}
			}

			assert.Equal(t, len(tt.wantMatch.Players), len(tt.args.match.Players))

			for wantPlayer, wantKillValue := range tt.wantMatch.Kills {
				for gotPlayer, gotKillValue := range tt.args.match.Kills {
					assert.Equal(t, wantPlayer, gotPlayer)
					assert.Equal(t, wantKillValue, gotKillValue)
				}
			}
		})
	}
}
