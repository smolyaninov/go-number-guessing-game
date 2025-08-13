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

	fmt.Println("🎯 Welcome to the Number Guessing Game!")
	fmt.Println("I'm thinking of a number between 1 and 100.")
	fmt.Println("You have to guess it based on the difficulty you choose.")
	fmt.Println()

	for {
		difficulty := selectDifficulty()
		chances := getChancesByDifficulty(difficulty)
		secret := rand.Intn(100) + 1

		if hs, err := hsService.Get(); err == nil {
			e := hs[difficulty]
			if e.Attempts > 0 {
				fmt.Printf("🏆 Current high score for %s: %d attempt(s), %.2fs (%s)\n",
					difficulty, e.Attempts, e.DurationSeconds, e.AchievedAt.Format(time.RFC3339))
			} else {
				fmt.Printf("🏆 No high score yet for %s. Be the first!\n", difficulty)
			}
		} else {
			fmt.Printf("Could not load high scores: %v\n", err)
		}

		fmt.Printf("\nGreat! You selected %s. You have %d chances.\n", difficulty, chances)

		startTime := time.Now()
		win, attempts := playGame(secret, chances)
		duration := time.Since(startTime).Seconds()

		if win {
			fmt.Println("🎉 Congratulations! You guessed the correct number!")

			if updated, err := hsService.UpdateIfBetter(difficulty, attempts, duration); err != nil {
				fmt.Printf("Failed to update high score: %v\n", err)
			} else if updated {
				fmt.Printf("🥇 New high score for %s: %d attempt(s), %.2fs!\n", difficulty, attempts, duration)
			}
		} else {
			fmt.Printf("💀 You've run out of chances. The number was %d.\n", secret)
		}

		fmt.Printf("🕐 You spent %.2f seconds.\n", duration)

		if hs, err := hsService.Get(); err == nil {
			printHighScores(hs)
		} else {
			fmt.Printf("Could not load high scores: %v\n", err)
		}

		var again string
		fmt.Print("\nDo you want to play again? (y/n): ")
		fmt.Scanln(&again)
		if again != "y" && again != "Y" {
			fmt.Println("Thanks for playing!")
			break
		}
	}
}

func selectDifficulty() domain.Level {
	for {
		fmt.Println("Select difficulty:")
		fmt.Println("1. Easy (10 chances)")
		fmt.Println("2. Medium (5 chances)")
		fmt.Println("3. Hard (3 chances)")
		fmt.Print("Enter choice (1/2/3): ")

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
			fmt.Println("Invalid choice. Please try again.")
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
	fmt.Println("💡 You have 1 hint available! Enter -1 to use it.")

	for i := 1; i <= chances; i++ {
		var guess int
		fmt.Printf("Attempt %d: Enter your guess: ", i)
		fmt.Scanln(&guess)

		if guess == -1 {
			if !hintUsed {
				printHint(secret)
				hintUsed = true
			} else {
				fmt.Println("You already used your hint!")
			}
			continue
		}

		if guess == secret {
			fmt.Printf("✅ Correct! You guessed it in %d attempts.\n", i)
			return true, i
		} else if guess < secret {
			fmt.Println("🔺 Too low. Try a higher number.")
		} else {
			fmt.Println("🔻 Too high. Try a lower number.")
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

	fmt.Println("💡 Hint:")
	fmt.Printf("👉 The number is between %d and %d.\n", low, high)

	if secret%2 == 0 {
		fmt.Println("⚖️ It's an even number.")
	} else {
		fmt.Println("⚖️ It's an odd number.")
	}
}

func printHighScores(hs map[domain.Level]domain.HighScore) {
	fmt.Println("\n===== High Scores =====")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "LEVEL\tATTEMPTS\tDURATION(s)\tDATE")

	for _, lvl := range domain.AllLevels() {
		e := hs[lvl]
		if e.Attempts > 0 {
			fmt.Fprintf(
				w,
				"%s\t%d\t%.2f\t%s\n",
				lvl,
				e.Attempts,
				e.DurationSeconds,
				e.AchievedAt.Format("2006-01-02 15:04:05"),
			)
		} else {
			fmt.Fprintf(w, "%s\t—\t—\t—\n", lvl)
		}
	}

	w.Flush()
	fmt.Println("=======================")
}
