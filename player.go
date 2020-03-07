package main

import (
	"fmt"
	"time"
)

// Ship is a ship.
type Ship struct {
	In         bool
	LastWarpIn time.Time
	LastDeath  time.Time
}

func NewShip(in bool, lwi, ld time.Time) Ship {
	return Ship{in, lwi, ld}
}

func (s *Ship) Destroy(timeOfDestruction time.Time) {
	s.In = false
	s.LastDeath = timeOfDestruction
}

func (s *Ship) DestroyedSince(since string) error {
	dur, err := time.ParseDuration(since)
	if err == nil {
		s.Destroy(time.Now().Add(-dur))
	}
	return err
}

func (s *Ship) WarpIn(timeOfWarpIn time.Time) {
	s.In = true
	s.LastWarpIn = timeOfWarpIn
}

func (s *Ship) WarpedInSince(since string) error {
	dur, err := time.ParseDuration(since)
	if err == nil {
		s.WarpIn(time.Now().Add(-dur))
	}
	return err
}

func (s *Ship) WarpOut() {
	s.In = false
}

func (s *Ship) Status(clock Clock) string {
	cooldownDuration := Cooldown(clock, s.LastWarpIn, s.LastDeath)
	statusLabel := "Out"
	if s.In {
		statusLabel = "In"
	}
	roundDuration := cooldownDuration.Round(time.Second)
	cooldownLabel := roundDuration.String()
	if !s.In && roundDuration <= 0 {
		cooldownLabel = "READY"
	}
	return "[" + statusLabel + "][" + cooldownLabel + "]"
}

// Player is a guy who has ships
type Player struct {
	Name                string
	Level               int
	Battleship, Support *Ship
}

func NewPlayer(name string, level int) Player {
	zeroTime := time.Time{}
	return Player{name, level, &Ship{false, zeroTime, zeroTime}, &Ship{false, zeroTime, zeroTime}}
}

func (p *Player) Status(clock Clock) string {
	return fmt.Sprintf("%s[%d] Battleship%s Support%s",
		p.Name, p.Level, p.Battleship.Status(clock), p.Support.Status(clock))
}

// func (p *Player) GetStatus() string {
// 	nameAndLevel := p.Name + "[" + string(p.Level) + "]"
// 	bsInOut := "Out"
// 	if p.Battleship.In {
// 		bsInOut = "In"
// 	}
// 	supportInOut := "Out"
// 	if p.Support.In {
// 		supportInOut = "In"
// 	}
// 	return ""
// }
