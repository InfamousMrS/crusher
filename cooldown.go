package main

import "time"

const (
	WarpInCooldownDuration time.Duration = time.Hour * 2
	DeathCooldownDuration  time.Duration = time.Hour * 18
)

type Clock interface {
	Now() time.Time
}

type SystemClock struct{}

func (s *SystemClock) Now() time.Time {
	return time.Now()
}

func Cooldown(clock Clock, lastWarpIn, lastDeath time.Time) time.Duration {
	now := clock.Now()
	timeSinceLastWarpIn := now.Sub(lastWarpIn)
	timeSinceLastDeath := now.Sub(lastDeath)
	remainingCooldown := DeathCooldownDuration - timeSinceLastDeath
	if remainingCooldown > 0 {
		return remainingCooldown
	} else {
		remainingCooldown = WarpInCooldownDuration - timeSinceLastWarpIn
		if remainingCooldown > 0 {
			return remainingCooldown
		}
	}
	return 0
}
