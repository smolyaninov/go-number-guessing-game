package service

import (
	"github.com/smolyaninov/go-number-guessing-game/internal/domain"
	"github.com/smolyaninov/go-number-guessing-game/internal/repo"
	"time"
)

type HighScoreService struct {
	repo repo.HighScoreRepository
}

func NewHighScoreService(repo repo.HighScoreRepository) *HighScoreService {
	return &HighScoreService{repo: repo}
}

func (s *HighScoreService) Get() (domain.HighScores, error) {
	hs, err := s.repo.Load()
	if err != nil {
		return nil, err
	}
	if hs == nil {
		hs = domain.HighScores{
			domain.LevelEasy:   {Level: domain.LevelEasy},
			domain.LevelMedium: {Level: domain.LevelMedium},
			domain.LevelHard:   {Level: domain.LevelHard},
		}
	}
	return hs, nil
}

func (s *HighScoreService) UpdateIfBetter(level domain.Level, attempts int, durationSeconds float64) (bool, error) {
	hs, err := s.Get()
	if err != nil {
		return false, err
	}

	cur := hs[level]
	shouldUpdate := cur.Attempts == 0 ||
		attempts < cur.Attempts ||
		(attempts == cur.Attempts && durationSeconds < cur.DurationSeconds)

	if !shouldUpdate {
		return false, nil
	}

	hs[level] = domain.HighScore{
		Level:           level,
		Attempts:        attempts,
		DurationSeconds: durationSeconds,
		AchievedAt:      time.Now().UTC(),
	}

	if err := s.repo.Save(hs); err != nil {
		return false, err
	}
	return true, nil
}
