package domain

type Level string

const (
	LevelEasy   Level = "Easy"
	LevelMedium Level = "Medium"
	LevelHard   Level = "Hard"
)

func AllLevels() []Level {
	return []Level{LevelEasy, LevelMedium, LevelHard}
}
