package domain

type Level string

const (
	LevelEasy   Level = "Easy"
	LevelMedium Level = "Medium"
	LevelHard   Level = "Hard"
)

const (
	ChancesEasy   = 10
	ChancesMedium = 5
	ChancesHard   = 3
)

var levelChances = map[Level]int{
	LevelEasy:   ChancesEasy,
	LevelMedium: ChancesMedium,
	LevelHard:   ChancesHard,
}

func AllLevels() []Level {
	return []Level{LevelEasy, LevelMedium, LevelHard}
}

func ChancesByLevel(level Level) int {
	if v, ok := levelChances[level]; ok {
		return v
	}
	return ChancesMedium
}
