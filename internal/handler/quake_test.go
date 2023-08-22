package handler

import (
	"bufio"
	"reflect"
	"strings"
	"testing"

	"github.com/maaarkin/quakecrawler/internal/domain"
)

func Test_quakeHandler_getPlayerName(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "with valid input player name",
			args: args{
				line: `21:15 ClientUserinfoChanged: 2 n\Isgalamido\t\0\model\uriel/zael\hmodel\uriel/zael\g_redteam\\g_blueteam\\c1\5\c2\5\hc\100\w\0\l\0\tt\0\tl\0`,
			},
			want: "Isgalamido",
		},
		{
			name: "with invalid input player name",
			args: args{
				line: `n\Isgalamido\t\0\model\uriel/zael\hmodel\uriel/zael\g_redteam\\g_blueteam\\c1\5\c2\5\hc\100\w\0\l\0\tt\0\tl\0`,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := NewQuakeHandler()
			if got := q.getPlayerName(tt.args.line); got != tt.want {
				t.Errorf("quakeHandler.getPlayerName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_quakeHandler_getDataFromKill(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name          string
		args          args
		wantKiller    string
		wantDeath     string
		wantDeathType string
	}{
		{
			name: "with valid input kill data",
			args: args{
				line: `21:42 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT`,
			},
			wantKiller:    "<world>",
			wantDeath:     "Isgalamido",
			wantDeathType: "MOD_TRIGGER_HURT",
		},
		{
			name: "with invalid input kill data",
			args: args{
				line: `1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT`,
			},
			wantKiller:    "",
			wantDeath:     "",
			wantDeathType: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := NewQuakeHandler()
			gotKiller, gotDeath, gotDeathType := q.getDataFromKill(tt.args.line)
			if gotKiller != tt.wantKiller {
				t.Errorf("quakeHandler.getDataFromKill() gotKiller = %v, want %v", gotKiller, tt.wantKiller)
			}
			if gotDeath != tt.wantDeath {
				t.Errorf("quakeHandler.getDataFromKill() gotDeath = %v, want %v", gotDeath, tt.wantDeath)
			}
			if gotDeathType != tt.wantDeathType {
				t.Errorf("quakeHandler.getDataFromKill() gotDeathType = %v, want %v", gotDeathType, tt.wantDeathType)
			}
		})
	}
}

func Test_quakeHandler_Run(t *testing.T) {
	type args struct {
		scanner *bufio.Scanner
	}
	tests := []struct {
		name          string
		args          func() args
		wantReport    map[string]domain.Payload
		wantKBMReport map[string]map[string]uint64
	}{
		{
			name: "with valid log payload",
			args: func() args {
				const input = `
				0:00 ------------------------------------------------------------
				0:00 InitGame: \sv_floodProtect\1\sv_maxPing\0\sv_minPing\0\sv_maxRate\10000\sv_minRate\0\sv_hostname\Code Miner Server\g_gametype\0\sv_privateClients\2\sv_maxclients\16\sv_allowDownload\0\dmflags\0\fraglimit\20\timelimit\15\g_maxGameClients\0\capturelimit\8\version\ioq3 1.36 linux-x86_64 Apr 12 2009\protocol\68\mapname\q3dm17\gamename\baseq3\g_needpass\0
				15:00 Exit: Timelimit hit.
				20:34 ClientConnect: 2
				20:34 ClientUserinfoChanged: 2 n\Isgalamido\t\0\model\xian/default\hmodel\xian/default\g_redteam\\g_blueteam\\c1\4\c2\5\hc\100\w\0\l\0\tt\0\tl\0
				20:37 ClientUserinfoChanged: 2 n\Isgalamido\t\0\model\uriel/zael\hmodel\uriel/zael\g_redteam\\g_blueteam\\c1\5\c2\5\hc\100\w\0\l\0\tt\0\tl\0
				20:37 ClientBegin: 2
				20:37 ShutdownGame:
				`
				scanner := bufio.NewScanner(strings.NewReader(input))

				return args{
					scanner: scanner,
				}
			},
			wantReport: map[string]domain.Payload{
				"game_1": {
					TotalKills: 0,
					Players:    []string{"Isgalamido"},
					Kills:      make(map[string]int64),
				},
			},
			wantKBMReport: make(map[string]map[string]uint64),
		},
		{
			name: "with valid log payload with kill",
			args: func() args {
				const input = `
				0:00 ------------------------------------------------------------
				0:00 InitGame: \sv_floodProtect\1\sv_maxPing\0\sv_minPing\0\sv_maxRate\10000\sv_minRate\0\sv_hostname\Code Miner Server\g_gametype\0\sv_privateClients\2\sv_maxclients\16\sv_allowDownload\0\dmflags\0\fraglimit\20\timelimit\15\g_maxGameClients\0\capturelimit\8\version\ioq3 1.36 linux-x86_64 Apr 12 2009\protocol\68\mapname\q3dm17\gamename\baseq3\g_needpass\0
				15:00 Exit: Timelimit hit.
				20:34 ClientConnect: 2
				20:34 ClientUserinfoChanged: 2 n\Isgalamido\t\0\model\xian/default\hmodel\xian/default\g_redteam\\g_blueteam\\c1\4\c2\5\hc\100\w\0\l\0\tt\0\tl\0
				20:37 ClientUserinfoChanged: 2 n\Isgalamido\t\0\model\uriel/zael\hmodel\uriel/zael\g_redteam\\g_blueteam\\c1\5\c2\5\hc\100\w\0\l\0\tt\0\tl\0
				20:37 ClientBegin: 2
				20:44 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT
				20:47 ShutdownGame:
				`
				scanner := bufio.NewScanner(strings.NewReader(input))

				return args{
					scanner: scanner,
				}
			},
			wantReport: map[string]domain.Payload{
				"game_1": {
					TotalKills: 1,
					Players:    []string{"Isgalamido"},
					Kills: map[string]int64{
						"Isgalamido": -1,
					},
				},
			},
			wantKBMReport: map[string]map[string]uint64{
				"game_1": {
					"MOD_TRIGGER_HURT": 1,
				},
			},
		},
		{
			name: "with valid log payload with multiples games",
			args: func() args {
				const input = `
				0:00 ------------------------------------------------------------
				0:00 InitGame: \sv_floodProtect\1\sv_maxPing\0\sv_minPing\0\sv_maxRate\10000\sv_minRate\0\sv_hostname\Code Miner Server\g_gametype\0\sv_privateClients\2\sv_maxclients\16\sv_allowDownload\0\dmflags\0\fraglimit\20\timelimit\15\g_maxGameClients\0\capturelimit\8\version\ioq3 1.36 linux-x86_64 Apr 12 2009\protocol\68\mapname\q3dm17\gamename\baseq3\g_needpass\0
				15:00 Exit: Timelimit hit.
				20:34 ClientConnect: 2
				20:34 ClientUserinfoChanged: 2 n\Isgalamido\t\0\model\xian/default\hmodel\xian/default\g_redteam\\g_blueteam\\c1\4\c2\5\hc\100\w\0\l\0\tt\0\tl\0
				20:37 ClientUserinfoChanged: 2 n\Isgalamido\t\0\model\uriel/zael\hmodel\uriel/zael\g_redteam\\g_blueteam\\c1\5\c2\5\hc\100\w\0\l\0\tt\0\tl\0
				20:37 ClientBegin: 2
				20:44 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT
				20:47 ShutdownGame:
				20:48 ------------------------------------------------------------
 				20:48 ------------------------------------------------------------
				20:51 InitGame: \sv_floodProtect\1\sv_maxPing\0\sv_minPing\0\sv_maxRate\10000\sv_minRate\0\sv_hostname\Code Miner Server\g_gametype\0\sv_privateClients\2\sv_maxclients\16\sv_allowDownload\0\bot_minplayers\0\dmflags\0\fraglimit\20\timelimit\15\g_maxGameClients\0\capturelimit\8\version\ioq3 1.36 linux-x86_64 Apr 12 2009\protocol\68\mapname\q3dm17\gamename\baseq3\g_needpass\0
				20:51 ClientConnect: 2
				20:52 ClientUserinfoChanged: 2 n\Isgalamido\t\0\model\uriel/zael\hmodel\uriel/zael\g_redteam\\g_blueteam\\c1\5\c2\5\hc\100\w\0\l\0\tt\0\tl\0
				20:52 ClientBegin: 2
				20:53 Item: 2 weapon_rocketlauncher
				20:53 Item: 2 ammo_rockets
				20:54 Item: 2 item_armor_body
				20:55 ShutdownGame:
				`
				scanner := bufio.NewScanner(strings.NewReader(input))

				return args{
					scanner: scanner,
				}
			},
			wantReport: map[string]domain.Payload{
				"game_1": {
					TotalKills: 1,
					Players:    []string{"Isgalamido"},
					Kills: map[string]int64{
						"Isgalamido": -1,
					},
				},
				"game_2": {
					TotalKills: 0,
					Players:    []string{"Isgalamido"},
					Kills:      make(map[string]int64),
				},
			},
			wantKBMReport: map[string]map[string]uint64{
				"game_1": {
					"MOD_TRIGGER_HURT": 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &quakeHandler{}
			scanner := tt.args().scanner
			got, got1 := q.Run(scanner)
			if !reflect.DeepEqual(got, tt.wantReport) {
				t.Errorf("quakeHandler.Run() got = %v, want %v", got, tt.wantReport)
			}
			if !reflect.DeepEqual(got1, tt.wantKBMReport) {
				t.Errorf("quakeHandler.Run() got1 = %v, want %v", got1, tt.wantKBMReport)
			}
		})
	}
}
