# Number Guessing Game CLI

This is my fourth project in Go, completed as part of
the [Number Guessing Game project](https://roadmap.sh/projects/number-guessing-game) on roadmap.sh.

The goal was to practice Go basics with user input, control flow, modular project structure, and simple persistence — by
implementing a classic terminal game with scoring and difficulty levels.

The application is a command-line game where you try to guess a random number within limited chances, optionally using a
one-time hint, and compete against your own high scores.

## Features

- Three difficulty levels:
    - Easy (10 chances)
    - Medium (5 chances)
    - Hard (3 chances)
- Random secret number between 1 and 100
- User feedback on each guess (“Higher” / “Lower”)
- One-time **hint system**: shows narrowed range + odd/even parity
- **Timer**: shows how long each round took
- **High scores**: per difficulty, saved locally in JSON
- **Replay loop**: keep playing until you quit
- Clean architecture with domain, engine, repo, and service layers

## Usage

```bash
guessing-game              # start the game
# Follow the prompts:
# - select difficulty (1/2/3)
# - enter guesses
# - optionally type -1 once to use the hint
```

### Example session

```
🎯 Number Guessing Game — guess a number between 1 and 100.

Select difficulty:
    1. Easy (10 chances)
    2. Medium (5 chances)
    3. Hard (3 chances)

Enter choice (1/2/3): 2

Level: Medium • Chances: 5 • Range: 1..100

Hint available once: enter -1 to use it.

Attempt 1/5: 50
Lower

Attempt 2/5: 25
Higher

Attempt 3/5: 30
✅ Correct in 3 attempt(s). Time: 7.12s
🥇 High score updated.
```

## Build

```bash
go build -o guessing-game ./cmd/guessing-game/main.go
```
