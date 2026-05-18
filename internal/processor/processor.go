package processor

import (
	"dungeon-challenge/internal/config"
	"dungeon-challenge/internal/event"
	"dungeon-challenge/internal/player"
	"fmt"
	"strconv"
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
	if pl.Finished {
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
			pl.BossEntered ||
			pl.CompletedFloors[pl.CurrentFloor] ||
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
	// next floor
	case 4:
		if !pl.InDungeon ||
			pl.BossEntered ||
			!pl.CompletedFloors[pl.CurrentFloor] ||
			pl.CurrentFloor >= p.config.Floors {
			p.impossibleMove(event)

			return
		}
		pl.CurrentFloor++
		if !pl.CompletedFloors[pl.CurrentFloor] {
			pl.FloorStartTime = event.Time
		}
		fmt.Printf(
			"[%s] Player [%d] went to the next floor\n",
			formatTime(event.Time),
			event.PlayerID,
		)
	// prev floor
	case 5:
		if !pl.InDungeon ||
			pl.BossEntered ||
			pl.CurrentFloor == 1 {
			p.impossibleMove(event)
			return
		}
		pl.CurrentFloor--
		fmt.Printf(
			"[%s] Player [%d] went to the previous floor\n",
			formatTime(event.Time),
			event.PlayerID,
		)
	//entered the boss's floor
	case 6:
		if !pl.InDungeon ||
			len(pl.CompletedFloors) != p.config.Floors ||
			pl.BossEntered {
			p.impossibleMove(event)
			return
		}
		pl.BossEntered = true
		fmt.Printf(
			"[%s] Player [%d] entered the boss's floor\n",
			formatTime(event.Time),
			event.PlayerID,
		)
		//killed the boss
	case 7:
		if !pl.InDungeon ||
			!pl.BossEntered ||
			pl.BossKilled {
			p.impossibleMove(event)
			return
		}
		pl.BossKilled = true
		pl.BossKillTime = event.Time
		fmt.Printf(
			"[%s] Player [%d] killed the boss\n",
			formatTime(event.Time),
			event.PlayerID,
		)
		//left the dungeon
	case 8:
		if !pl.InDungeon {
			p.impossibleMove(event)
			return
		}
		pl.InDungeon = false
		pl.Finished = true
		pl.ExitTime = event.Time
		fmt.Printf(
			"[%s] Player [%d] left the dungeon\n",
			formatTime(event.Time),
			event.PlayerID,
		)
		//cannot continue due to [`reason`]
	case 9:
		if !pl.InDungeon {
			p.impossibleMove(event)
			return
		}
		fmt.Printf(
			"[%s] Player [%d] cannot continue due to [%s]\n",
			formatTime(event.Time),
			event.PlayerID,
			event.Extra,
		)
		pl.Disqualified = true
		pl.Finished = true
	//has restored [`health`] of health
	case 10:
		if !pl.InDungeon {
			p.impossibleMove(event)
			return
		}
		fmt.Printf(
			"[%s] Player [%d] has restored [%s] of health\n",
			formatTime(event.Time),
			event.PlayerID,
			event.Extra,
		)
		heal, _ := strconv.Atoi(event.Extra)
		pl.HP += heal
		if pl.HP > 100 {
			pl.HP = 100
		}

	//recieved [`damage`] of damage
	case 11:
		if !pl.InDungeon {
			p.impossibleMove(event)
			return
		}
		fmt.Printf(
			"[%s] Player [%d] recieved [%s] of damage\n",
			formatTime(event.Time),
			event.PlayerID,
			event.Extra,
		)
		dmg, _ := strconv.Atoi(event.Extra)
		pl.HP -= dmg
		if pl.HP <= 0 {
			pl.HP = 0
			pl.Dead = true
			fmt.Printf(
				"[%s] Player [%d] is dead\n",
				formatTime(event.Time),
				event.PlayerID,
			)
			pl.Finished = true
		}
	}
}
