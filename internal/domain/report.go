package domain

import "strings"

type Payload struct {
	TotalKills int64            `json:"total_kills"`
	Players    []string         `json:"players"`
	Kills      map[string]int64 `json:"kills,omitempty"`
}

func NewPayload() Payload {
	return Payload{
		TotalKills: 0,
		Players:    []string{},
		Kills:      make(map[string]int64),
	}
}

func (p *Payload) AddPlayer(name string) {
	for _, player := range p.Players {
		if player == name {
			return
		}
	}

	p.Players = append(p.Players, name)
}

func (p *Payload) AddKill(killer, died string) {
	p.TotalKills++

	if killer == died {
		return
	}

	if strings.Contains(killer, "world") {
		p.Kills[died] = p.Kills[died] - 1
	} else {
		p.Kills[killer] = p.Kills[killer] + 1
	}
}
