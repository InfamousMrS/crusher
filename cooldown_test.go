package main

import (
	"testing"
	"time"
)

type MockClock struct {
	NowTime time.Time
}

func (m MockClock) Now() time.Time {
	return m.NowTime
}

func TestCooldown(t *testing.T) {

	var (
		now                         = time.Now()
		fortyFiveMinutesAgo         = now.Add(-time.Minute * 45)
		twoHoursAgo                 = now.Add(-time.Hour * 2)
		twoHoursAndOneMinuteAgo     = now.Add(-time.Minute * 121)
		elevenHoursElevenMinutesAgo = now.Add((time.Hour + time.Minute) * -11)
		seventeenHoursAgo           = now.Add(-time.Hour * 17)
		eighteenHoursAgo            = now.Add(-time.Hour * 18)
		foreverAgo                  = now.Add(-time.Hour * 10000)
	)

	t.Run("Cooldown Right After Warp In Times", func(t *testing.T) {
		mockClock := MockClock{now}
		warpInTimes := []struct {
			lastWarpIn   time.Time
			lastDeath    time.Time
			wantCooldown time.Duration
		}{
			{lastWarpIn: now, lastDeath: foreverAgo, wantCooldown: WarpInCooldownDuration},
			{lastWarpIn: fortyFiveMinutesAgo, lastDeath: foreverAgo, wantCooldown: time.Minute * 75},
			{lastWarpIn: twoHoursAgo, lastDeath: foreverAgo, wantCooldown: 0},
			{lastWarpIn: twoHoursAndOneMinuteAgo, lastDeath: foreverAgo, wantCooldown: 0},
			{lastWarpIn: twoHoursAndOneMinuteAgo, lastDeath: now, wantCooldown: time.Hour * 18},
			{lastWarpIn: eighteenHoursAgo, lastDeath: elevenHoursElevenMinutesAgo, wantCooldown: time.Minute * 409},
			{lastWarpIn: eighteenHoursAgo, lastDeath: seventeenHoursAgo, wantCooldown: time.Hour},
			{lastWarpIn: eighteenHoursAgo, lastDeath: eighteenHoursAgo, wantCooldown: 0},
		}

		for _, testCase := range warpInTimes {
			got := Cooldown(mockClock, testCase.lastWarpIn, testCase.lastDeath)
			assertDurations(t, testCase.wantCooldown, got)
		}
	})

}

func assertDurations(t *testing.T, want, got time.Duration) {
	t.Helper()
	if got != want {
		t.Errorf("Expected %s but got %s", want, got)
	}
}
