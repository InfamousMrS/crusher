package main

import "testing"

func TestNewWsMatch(t *testing.T) {
	name := "Blah"
	wsmatch := NewWsMatch(name)

	if wsmatch.Name != name {
		t.Errorf("Expected match name %s but got %s", name, wsmatch.Name)
	}
	expectedLength := 0
	if len(wsmatch.Friendlies) != expectedLength {
		t.Errorf("Expected %d friendlies, but got %d", expectedLength, len(wsmatch.Friendlies))
	}
	if len(wsmatch.Enemies) != expectedLength {
		t.Errorf("Expected %d enemies, but got %d", expectedLength, len(wsmatch.Enemies))
	}
}

func TestAddEnemy(t *testing.T) {
	t.Run("Adding a player not yet added", func(t *testing.T) {
		enemy := NewPlayer("BlackDeath", 500)
		match := NewWsMatch("SystemShock")

		match.AddEnemy(enemy)
		if len(match.Enemies) != 1 {
			t.Errorf("Expected one enemy but got %d", len(match.Enemies))
		}
	})
	t.Run("Adding a player already added", func(t *testing.T) {
		enemy := NewPlayer("BlackDeath", 500)
		friend := NewPlayer("BlackDeath", 222)
		match := NewWsMatch("SystemShock")

		match.AddFriendly(friend)
		err := match.AddEnemy(enemy)
		if err == nil {
			t.Error("Expected an error but got none.")
		}
	})
}

func TestAddFriendly(t *testing.T) {
	friendly := NewPlayer("Miahwt", 500)
	match := NewWsMatch("SystemShock")

	match.AddFriendly(friendly)
	if len(match.Friendlies) != 1 {
		t.Errorf("Expected one friend but got %d", len(match.Friendlies))
	}
}

// func TestReportAvailable(t *testing.T) {
// 	match := NewWsMatch("TestMatch")
// 	match.AddEnemy(NewPlayer("BlackDeath", 500))
// 	match.AddEnemy(NewPlayer("Singularity", 240))
// 	match.AddFriendly(NewPlayer("Hades", 222))
// 	match.AddFriendly(NewPlayer("Vicious", 100))

// 	report := match.ReportAvailable()
// 	expected := "```\nEnemies:\n   BlackDeath\n   Singularity\n\nFriendlies:\n   Hades\n   Vicious\n```"

// 	if expected != report {
// 		t.Errorf("Expected report [%s] but got [%s]", expected, report)
// 	}
// }

// func TestGetPlayerStatus(t *testing.T) {
// 	match := NewWsMatch("TestMatch")
// 	enemy := NewPlayer("BlackDeath", 233)
// 	friend := NewPlayer("Hades", 212)
// 	match.AddEnemy(enemy)
// 	match.AddFriendly(friend)
// 	enemy.Battleship.DestroyedSince("12h")
// 	enemy.Support.WarpedInSince("0h0m")
// 	friend.Support.WarpedInSince("1h5m")

// 	t.Run("Testing enemy status", func(t *testing.T) {
// 		got := match.GetPlayerStatus("BlackDeath")
// 		want := "```\nEnemy[BlackDeath]\n  Battleship[Out] Cooldown[6h0m]\n  Support[Out] Cooldown[2h0m]\n```"
// 		if got != want {
// 			t.Errorf("Wanted %s but got %s", want, got)
// 		}
// 	})

// }
