package game

import "math/rand"

type Engine struct {
	Min, Max int
	Secret   int
	Chances  int
}

func NewEngine(min int, max int, secret int, chances int) *Engine {
	return &Engine{Min: min, Max: max, Secret: secret, Chances: chances}
}

func (e *Engine) InRange(n int) bool {
	return n >= e.Min && n <= e.Max
}

func (e *Engine) Compare(guess int) int {
	switch {
	case guess < e.Secret:
		return -1
	case guess > e.Secret:
		return 1
	default:
		return 0
	}
}

func (e *Engine) HintRange() (low, high int) {
	span := rand.Intn(9) + 6        // [6..14]
	offset := rand.Intn(span-1) + 1 // [1..span-1]

	low = e.Secret - offset
	high = low + span

	if low < e.Min {
		shift := e.Min - low
		low += shift
		high += shift
	}
	if high > e.Max {
		shift := high - e.Max
		low -= shift
		high -= shift
	}
	if low > e.Secret {
		low = e.Secret
	}
	if high < e.Secret {
		high = e.Secret
	}
	return
}

func (e *Engine) Hint() (low int, high int, party string) {
	low, high = e.HintRange()
	if e.Secret%2 == 0 {
		party = "even"
	} else {
		party = "odd"
	}
	return
}
