package processor

import (
	"dungeon-challenge/internal/config"
	"dungeon-challenge/internal/event"
	"dungeon-challenge/internal/player"
	"fmt"
	"time"
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

func formatTime(t time.Time) string {
	return t.Format("15:04:05")
}

func (p *Processor) impossibleMove(event event.Event) {
	fmt.Printf(
		"[%s] Player [%d] makes imposible move [%d]\n",
		formatTime(event.Time),
		event.PlayerID,
		event.EventID,
	)
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
			formatTime(event.Time),
			event.PlayerID,
		)
		return
	}

	switch event.EventID {
	//registration
	case 1:
		if pl.Registered {
			p.impossibleMove(event)
			return
		}
		pl.Registered = true
		fmt.Printf(
			"[%s] Player [%d] registered\n",
			formatTime(event.Time),
			event.PlayerID,
		)
	//enter
	case 2:
		if pl.InDungeon {
			p.impossibleMove(event)
			return
		}
		pl.InDungeon = true
		pl.CurrentFloor = 1
		pl.FloorStartTime = event.Time
		pl.EnterTime = event.Time
		fmt.Printf(
			"[%s] Player [%d] entered the dungeon\n",
			formatTime(event.Time),
			event.PlayerID,
		)
	//kill monster
	case 3:
		if !pl.InDungeon ||
			pl.CompletedFloors[pl.CurrentFloor] ||
			pl.BossEntered ||
			pl.FloorKills[pl.CurrentFloor] >= p.config.Monsters {
			p.impossibleMove(event)
			return
		}

		pl.FloorKills[pl.CurrentFloor]++
		fmt.Printf(
			"[%s] Player [%d] killed the monster\n",
			formatTime(event.Time),
			event.PlayerID,
		)
		if pl.FloorKills[pl.CurrentFloor] == p.config.Monsters {
			pl.CompletedFloors[pl.CurrentFloor] = true
			duration := event.Time.Sub(pl.FloorStartTime)
			pl.FloorDurations = append(pl.FloorDurations, duration)

		}
	}
}
