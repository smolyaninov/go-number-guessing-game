package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/smolyaninov/go-number-guessing-game/internal/domain"
	"github.com/smolyaninov/go-number-guessing-game/internal/game"
	"github.com/smolyaninov/go-number-guessing-game/internal/input"
	"github.com/smolyaninov/go-number-guessing-game/internal/repo"
	"github.com/smolyaninov/go-number-guessing-game/internal/service"
)

const (
	rangeMin   = 1
	rangeMax   = 100
	timeFormat = "2006-01-02 15:04"
)

func main() {
	repository := repo.NewJSONHighScoreRepository("data/highscores.json")
	hsService := service.NewHighScoreService(repository)

	fmt.Printf("🎯 Number Guessing Game — guess a number between %d and %d.\n", rangeMin, rangeMax)
	fmt.Println()

	for {
		level := selectDifficulty()
		chances := domain.ChancesByLevel(level)
		secret := rand.Intn(rangeMax-rangeMin+1) + rangeMin

		if hs, err := hsService.Get(); err == nil {
			if e := hs[level]; e.Attempts > 0 {
				fmt.Printf("\n🏆 %s high score: %d attempts, %.2fs (%s)\n",
					level, e.Attempts, e.DurationSeconds, e.AchievedAt.Format(timeFormat))
			}
		}

		fmt.Printf("\nLevel: %s • Chances: %d • Range: %d..%d\n\n", level, chances, rangeMin, rangeMax)

		engine := game.NewEngine(rangeMin, rangeMax, secret, chances)
		start := time.Now()
		win, attempts := playGame(engine)
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

		ans, err := input.ReadString("Play again? (y/n): ")
		if err != nil {
			fmt.Println("Invalid input, please enter Y or N.")
			fmt.Println()
			continue
		}
		if !strings.EqualFold(ans, "y") && !strings.EqualFold(ans, "yes") {
			fmt.Println("\nBye!")
			return
		}
		fmt.Println()
	}
}

func selectDifficulty() domain.Level {
	for {
		fmt.Println("Select difficulty:")
		fmt.Printf("\t1. Easy (%d chances)\n", domain.ChancesEasy)
		fmt.Printf("\t2. Medium (%d chances)\n", domain.ChancesMedium)
		fmt.Printf("\t3. Hard (%d chances)\n", domain.ChancesHard)

		choice, err := input.ReadInt("\nEnter choice (1/2/3): ")
		if err != nil {
			fmt.Println("Invalid input, please enter 1, 2 or 3.")
			fmt.Println()
			continue
		}

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

func playGame(e *game.Engine) (bool, int) {
	hintUsed := false

	fmt.Println("Hint available once: enter -1 to use it.")
	fmt.Println()

	for i := 1; i <= e.Chances; i++ {
		guess, err := input.ReadInt(fmt.Sprintf("Attempt %d/%d: ", i, e.Chances))
		if err != nil {
			fmt.Println("Invalid input, please enter a number.")
			fmt.Println()
			i--
			continue
		}

		if guess == -1 {
			if !hintUsed {
				low, high := e.HintRange()
				fmt.Printf("\n💡 Hint: %d..%d; %s\n\n", low, high, game.Parity(e.Secret))
				hintUsed = true
			} else {
				fmt.Println("Hint already used.")
				fmt.Println()
				i--
			}
			continue
		}

		if !e.InRange(guess) {
			fmt.Printf("Enter a number between %d and %d.\n\n", e.Min, e.Max)
			i--
			continue
		}

		switch e.Compare(guess) {
		case 0:
			return true, i

		case -1:
			fmt.Println("Higher")
			fmt.Println()
		case 1:
			fmt.Println("Lower")
			fmt.Println()
		}
	}
	return false, e.Chances
}

func printHighScores(hs map[domain.Level]domain.HighScore) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "LEVEL\tATTEMPTS\tDURATION(s)\tDATE")
	for _, lvl := range domain.AllLevels() {
		e := hs[lvl]
		if e.Attempts > 0 {
			fmt.Fprintf(w, "%s\t%d\t%.2f\t%s\n",
				lvl, e.Attempts, e.DurationSeconds, e.AchievedAt.Format(timeFormat))
		} else {
			fmt.Fprintf(w, "%s\t—\t—\t—\n", lvl)
		}
	}
	w.Flush()
}
