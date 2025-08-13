package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("🎯 Welcome to the Number Guessing Game!")
	fmt.Println("I'm thinking of a number between 1 and 100.")
	fmt.Println("You have to guess it based on the difficulty you choose.")
	fmt.Println()

	for {
		difficulty := selectDifficulty()
		chances := getChancesByDifficulty(difficulty)
		secret := rand.Intn(100) + 1

		fmt.Printf("\nGreat! You selected %s. You have %d chances.\n", difficulty, chances)

		startTime := time.Now()

		win := playGame(secret, chances)

		duration := time.Since(startTime).Seconds()

		if win {
			fmt.Println("🎉 Congratulations! You guessed the correct number!")
		} else {
			fmt.Printf("💀 You've run out of chances. The number was %d.\n", secret)
		}

		fmt.Printf("🕐 You spent %.2f seconds.\n", duration)

		var again string
		fmt.Print("\nDo you want to play again? (y/n): ")
		fmt.Scanln(&again)
		if again != "y" && again != "Y" {
			fmt.Println("Thanks for playing!")
			break
		}
	}
}

func selectDifficulty() string {
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
			return "Easy"
		case 2:
			return "Medium"
		case 3:
			return "Hard"
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func getChancesByDifficulty(level string) int {
	switch level {
	case "Easy":
		return 10
	case "Medium":
		return 5
	case "Hard":
		return 3
	default:
		return 5
	}
}

func playGame(secret, chances int) bool {
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
			i--
			continue
		}

		if guess == secret {
			fmt.Printf("✅ Correct! You guessed it in %d attempts.\n", i)
			return true
		} else if guess < secret {
			fmt.Println("🔺 Too low. Try a higher number.")
		} else {
			fmt.Println("🔻 Too high. Try a lower number.")
		}
	}
	return false
}

func printHint(secret int) {
	const minVal, maxVal = 1, 100

	span := rand.Intn(9) + 6 // [6..14]

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
