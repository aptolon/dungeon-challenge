package processor

import (
	"dungeon-challenge/internal/config"
	"dungeon-challenge/internal/event"
	"dungeon-challenge/internal/player"
	"fmt"
)

type Processor struct {
	config  config.Config
	players map[int]*player.Player
}

func NewProcessor(cfg config.Config) *Processor {
	return &Processor{
		config:  cfg,
		players: make(map[int]*player.Player),
	}
}

func (p *Processor) getOrCreatePlayer(id int) *player.Player {
	existing, ok := p.players[id]
	if ok {
		return existing
	}
	newPlayer := player.NewPlayer(id)
	p.players[id] = newPlayer
	return newPlayer
}

func (p *Processor) Process(event event.Event) {
	pl := p.getOrCreatePlayer(event.PlayerID)
	if pl.Dead {
		return

	}
	if pl.Disqualified {
		return
	}
	if event.EventID != 1 && !pl.Registered {
		pl.Disqualified = true
		fmt.Printf(
			"[%s] Player [%d] is disqualified\n",
			event.Time.Format("15:04:05"),
			event.PlayerID,
		)
		return
	}

	switch event.EventID {
	case 1:
		pl.Registered = true
		fmt.Printf(
			"[%s] Player [%d] registered\n",
			event.Time.Format("15:04:05"),
			event.PlayerID,
		)
	case 2:
		pl.InDungeon = true
		pl.EnterTime = event.Time
		fmt.Printf(
			"[%s] Player [%d] entered the dungeon\n",
			event.Time.Format("15:04:05"),
			event.PlayerID,
		)
	}
}
