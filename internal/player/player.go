package player

import "time"

type Player struct {
	ID int

	Registered   bool
	InDungeon    bool
	Dead         bool
	Disqualified bool

	HP int

	CurrentFloor    int
	FloorKills      map[int]int
	CompletedFloors map[int]bool

	BossEntered bool
	BossKilled  bool

	EnterTime time.Time
	ExitTime  time.Time

	FloorStartTime time.Time
	FloorDurations []time.Duration

	BossEnterTime time.Time
	BossKillTime  time.Time
}

func NewPlayer(id int) *Player {
	return &Player{
		ID: id,
		HP: 100,

		FloorKills:      make(map[int]int),
		CompletedFloors: make(map[int]bool),
	}
}
