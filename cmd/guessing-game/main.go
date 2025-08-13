package main

import (
	"fmt"
	"math/rand"
	"os"
	"text/tabwriter"
	"time"

	"github.com/smolyaninov/go-number-guessing-game/internal/domain"
	"github.com/smolyaninov/go-number-guessing-game/internal/repo"
	"github.com/smolyaninov/go-number-guessing-game/internal/service"
)

func main() {
	repository := repo.NewJSONHighScoreRepository("data/highscores.json")
	hsService := service.NewHighScoreService(repository)

	fmt.Println("🎯 Number Guessing Game — guess a number between 1 and 100.")
	fmt.Println()

	for {
		level := selectDifficulty()
		chances := getChancesByDifficulty(level)
		secret := rand.Intn(100) + 1

		if hs, err := hsService.Get(); err == nil {
			if e := hs[level]; e.Attempts > 0 {
				fmt.Printf("\n🏆 %s high score: %d attempts, %.2fs (%s)\n",
					level, e.Attempts, e.DurationSeconds, e.AchievedAt.Format("2006-01-02 15:04"))
			}
		}

		fmt.Printf("\nLevel: %s • Chances: %d\n\n", level, chances)

		start := time.Now()
		win, attempts := playGame(secret, chances)
		duration := time.Since(start).Seconds()

		fmt.Println()
		if win {
			fmt.Printf("✅ Correct in %d attempt(s). Time: %.2fs\n", attempts, duration)
			if updated, err := hsService.UpdateIfBetter(level, attempts, duration); err == nil && updated {
				fmt.Println("🥇 High score updated.")
			} else if err != nil {
				fmt.Printf("! High score save failed: %v\n", err)
			}
		} else {
			fmt.Printf("💀 Out of chances. Number was %d. Time: %.2fs\n", secret, duration)
		}

		if hs, err := hsService.Get(); err == nil {
			fmt.Println()
			printHighScores(hs)
		}

		fmt.Print("\nPlay again? (y/n): ")
		var again string
		fmt.Scanln(&again)
		if again != "y" && again != "Y" {
			fmt.Println("\nBye!")
			return
		}
		fmt.Println()
	}
}

func selectDifficulty() domain.Level {
	for {
		fmt.Println("Select difficulty:")
		fmt.Println("\t1. Easy (10 chances)")
		fmt.Println("\t2. Medium (5 chances)")
		fmt.Println("\t3. Hard (3 chances)")
		fmt.Print("\nEnter choice (1/2/3): ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			return domain.LevelEasy
		case 2:
			return domain.LevelMedium
		case 3:
			return domain.LevelHard
		default:
			fmt.Println("Invalid choice, try again.")
			fmt.Println()
		}
	}
}

func getChancesByDifficulty(level domain.Level) int {
	switch level {
	case domain.LevelEasy:
		return 10
	case domain.LevelMedium:
		return 5
	case domain.LevelHard:
		return 3
	default:
		return 5
	}
}

func playGame(secret, chances int) (bool, int) {
	hintUsed := false

	for i := 1; i <= chances; i++ {
		fmt.Printf("Attempt %d/%d — guess (-1 = hint): ", i, chances)
		var guess int
		fmt.Scanln(&guess)

		if guess == -1 {
			if !hintUsed {
				printHint(secret)
				hintUsed = true
			} else {
				fmt.Println("Hint already used.")
				fmt.Println()
			}
			continue
		}

		if guess == secret {
			return true, i
		}
		if guess < secret {
			fmt.Println("Higher")
			fmt.Println()
		} else {
			fmt.Println("Lower")
			fmt.Println()
		}
	}
	return false, chances
}

func printHint(secret int) {
	const minVal, maxVal = 1, 100
	span := rand.Intn(9) + 6        // [6..14]
	offset := rand.Intn(span-1) + 1 // [1..span-1]

	low := secret - offset
	high := low + span

	if low < minVal {
		shift := minVal - low
		low += shift
		high += shift
	}
	if high > maxVal {
		shift := high - maxVal
		low -= shift
		high -= shift
	}
	if low > secret {
		low = secret
	}
	if high < secret {
		high = secret
	}

	fmt.Printf("\n💡 Hint: %d..%d; %s\n\n", low, high, parity(secret))
}

func parity(n int) string {
	if n%2 == 0 {
		return "even"
	}
	return "odd"
}

func printHighScores(hs map[domain.Level]domain.HighScore) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "LEVEL\tATTEMPTS\tDURATION(s)\tDATE")
	for _, lvl := range domain.AllLevels() {
		e := hs[lvl]
		if e.Attempts > 0 {
			fmt.Fprintf(w, "%s\t%d\t%.2f\t%s\n",
				lvl, e.Attempts, e.DurationSeconds, e.AchievedAt.Format("2006-01-02 15:04"))
		} else {
			fmt.Fprintf(w, "%s\t—\t—\t—\n", lvl)
		}
	}
	w.Flush()
}
