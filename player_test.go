package main

import (
	"testing"
	"time"
)

func TestDestroy(t *testing.T) {
	now := time.Now()
	foreverAgo := now.Add(-time.Hour * 10000)

	ship := Ship{true, foreverAgo, foreverAgo}
	ship.Destroy(now)

	want := Ship{false, foreverAgo, now}
	assertShip(t, ship, want)
}

func TestDestroyedSince(t *testing.T) {
	foreverAgo := time.Now().Add(-time.Hour * 10000)

	t.Run("Valid input", func(t *testing.T) {
		ship := Ship{true, foreverAgo, foreverAgo}
		before := time.Now()
		err := ship.DestroyedSince("2h5m")
		after := time.Now()

		expectedBefore := before.Add(-time.Minute * 125)
		expectedAfter := after.Add(-time.Minute * 125)

		assertNoError(t, err)
		assertBetweenTimes(t, ship.LastDeath, expectedBefore, expectedAfter)
	})

	t.Run("Invalid duration string", func(t *testing.T) {
		ship := Ship{true, foreverAgo, foreverAgo}
		err := ship.DestroyedSince("a while back")

		assertShip(t, ship, Ship{true, foreverAgo, foreverAgo})
		assertError(t, err)
	})
}

func TestWarpIn(t *testing.T) {
	now := time.Now()
	foreverAgo := time.Now().Add(-time.Hour * 10000)

	ship := Ship{In: false, LastWarpIn: foreverAgo, LastDeath: foreverAgo}
	ship.WarpIn(now)

	want := Ship{In: true, LastWarpIn: now, LastDeath: foreverAgo}
	assertShip(t, ship, want)

}

func TestWarpedInSince(t *testing.T) {
	foreverAgo := time.Now().Add(-time.Hour * 10000)

	t.Run("Valid duration string", func(t *testing.T) {
		ship := Ship{false, foreverAgo, foreverAgo}
		before := time.Now()
		err := ship.WarpedInSince("4h20m")
		after := time.Now()

		expectedBefore := before.Add(-time.Minute * 260)
		expectedAfter := after.Add(-time.Minute * 260)

		assertNoError(t, err)
		assertBetweenTimes(t, ship.LastWarpIn, expectedBefore, expectedAfter)
	})

	t.Run("Invalid duration", func(t *testing.T) {
		ship := Ship{false, foreverAgo, foreverAgo}
		err := ship.WarpedInSince("hhmm")

		assertShip(t, ship, Ship{false, foreverAgo, foreverAgo})
		assertError(t, err)
	})
}

func TestWarpOut(t *testing.T) {
	aLittleWhileBack := time.Now().Add(-time.Minute * 20)
	foreverAgo := time.Now().Add(-time.Hour * 10000)

	ship := Ship{In: true, LastWarpIn: aLittleWhileBack, LastDeath: foreverAgo}
	ship.WarpOut()

	want := Ship{In: false, LastWarpIn: aLittleWhileBack, LastDeath: foreverAgo}
	assertShip(t, ship, want)
}

func TestNewPlayer(t *testing.T) {
	nulShip := Ship{false, time.Time{}, time.Time{}}
	player := NewPlayer("Bob", 202)

	assertPlayer(t, Player{"Bob", 202, nulShip, nulShip}, player)
}

func TestNewShip(t *testing.T) {
	now := time.Now()
	earlier := now.Add(-time.Hour)

	got := NewShip(false, now, earlier)
	want := Ship{false, now, earlier}

	if got != want {
		t.Errorf("Wanted %v but got %v", want, got)
	}
}

func TestShipStatus(t *testing.T) {

	now := time.Now()
	fiveMinsAgo := now.Add(-time.Minute * 5).Add(-time.Second)
	oneHourAgo := now.Add(-time.Hour)
	tminus1h59mAgo := now.Add(-time.Hour).Add(-time.Minute * 59)
	tminus1h59m59sAgo := tminus1h59mAgo.Add(-time.Second * 59)
	twoHourAgo := now.Add(-time.Hour * 2)
	zero := time.Time{}
	mockClock := MockClock{now}

	cases := []struct {
		In     bool
		WarpIn time.Time
		Death  time.Time
		Result string
	}{
		{In: true, WarpIn: now, Death: zero, Result: "[In][2h0m0s]"},
		{In: false, WarpIn: fiveMinsAgo, Death: zero, Result: "[Out][1h54m59s]"},
		{In: false, WarpIn: fiveMinsAgo, Death: oneHourAgo, Result: "[Out][17h0m0s]"},
		{In: true, WarpIn: tminus1h59mAgo, Death: zero, Result: "[In][1m0s]"},
		{In: true, WarpIn: tminus1h59m59sAgo, Death: zero, Result: "[In][1s]"},
		{In: true, WarpIn: twoHourAgo, Death: zero, Result: "[In][0s]"},
		{In: false, WarpIn: twoHourAgo, Death: zero, Result: "[Out][READY]"},
	}

	for _, tcase := range cases {
		ship := NewShip(tcase.In, tcase.WarpIn, tcase.Death)
		actual := ship.Status(mockClock)

		if actual != tcase.Result {
			t.Errorf("Expected %s but got %s for ship %v", tcase.Result, actual, ship)
		}
	}
}

func TestPlayerStatus(t *testing.T) {
	now := time.Now()
	zero := time.Time{}
	mockClock := MockClock{now}

	justWarpedInShip := NewShip(true, now, zero)
	justDiedShip := NewShip(false, zero, now)
	cleanShip := NewShip(true, zero, zero)
	cleanOutShip := NewShip(false, zero, zero)

	cases := []struct {
		P           Player
		TestBS      Ship
		TestSupport Ship
		Result      string
	}{
		{NewPlayer("Hades", 120), justDiedShip, cleanOutShip,
			"Hades[120] Battleship[Out][18h0m0s] Support[Out][READY]"},
		{NewPlayer("Miahwt", 99), justWarpedInShip, cleanShip,
			"Miahwt[99] Battleship[In][2h0m0s] Support[In][0s]"},
		{NewPlayer("Infamosu", 210), cleanOutShip, justDiedShip,
			"Infamosu[210] Battleship[Out][READY] Support[Out][18h0m0s]"},
	}

	for _, tcase := range cases {
		player := Player{tcase.P.Name, tcase.P.Level, tcase.TestBS, tcase.TestSupport}
		actual := player.Status(mockClock)

		if actual != tcase.Result {
			t.Errorf("Expected %s but got %s for player %v", tcase.Result, actual, player)
		}
	}
}

func assertPlayer(t *testing.T, want, got Player) {
	t.Helper()
	if want != got {
		t.Errorf("Expected player %v but got %v", want, got)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("Expecting no error but got %s", err)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Error("Expecting an error but got nil")
	}
}

func assertShip(t *testing.T, got, want Ship) {
	t.Helper()
	if want != got {
		t.Errorf("Wanted %v but got %v", want, got)
	}
}

func assertBetweenTimes(t *testing.T, got, before, after time.Time) {
	t.Helper()
	if got.Before(before) || after.Before(got) {
		t.Errorf("Expected a time between %s and %s, but got %s", before, after, got)
	}
}
