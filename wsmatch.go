package main

import (
	"fmt"
)

type WsMatch struct {
	Name       string
	Enemies    []Player
	Friendlies []Player
	players    map[string]*Player
}

func NewWsMatch(name string) WsMatch {
	friends := []Player{}
	enemies := []Player{}
	players := map[string]*Player{}
	return WsMatch{name, friends, enemies, players}
}

func (wsm *WsMatch) AddEnemy(p Player) error {
	return wsm.addPlayer(p, false)
}

func (wsm *WsMatch) AddFriendly(p Player) error {
	return wsm.addPlayer(p, true)
}

func (wsm *WsMatch) addPlayer(p Player, friendly bool) error {
	_, exists := wsm.players[p.Name]
	if !exists {
		wsm.players[p.Name] = &p
		if friendly {
			wsm.Friendlies = append(wsm.Friendlies, p)
		} else {
			wsm.Enemies = append(wsm.Enemies, p)
		}
	} else {
		return fmt.Errorf("Player %s is already in the Star", p.Name)
	}
	return nil
}

func (wsm *WsMatch) ReportAvailable() string {
	// report := "```\nEnemies\n"

	return ""
}

func (wsm *WsMatch) GetPlayerStatus(name string) string {
	return ""
}
