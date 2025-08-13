package domain

import "time"

type HighScore struct {
	Level           Level     `json:"level"`
	Attempts        int       `json:"attempts"`
	DurationSeconds float64   `json:"duration_seconds"`
	AchievedAt      time.Time `json:"achieved_at"`
}

type HighScores map[Level]HighScore
