package main

import (
	"fmt"
	"os"
	"strings"
)

func Fatal(err error, msg string) {
	if err != nil {
		panic(msg)
	}
}

type BattleResult int

const (
	Win BattleResult = iota
	Lose
	Draw
)

/*
	Part 1
*/

func DoBattle(a string, b string) BattleResult {
	switch a {
	case "A": // Rock
		switch b {
		case "Z": // Scissors
			return Lose
		case "Y": // Paper
			return Win
		}
	case "B": // Paper
		switch b {
		case "X": // Rock
			return Lose
		case "Z": // Paper
			return Win
		}
	case "C": // Scissors
		switch b {
		case "Y": // Paper
			return Lose
		case "X": // Rock
			return Win
		}
	default:
		panic("bad battle: " + a + " vs " + b)
	}
	return Draw
}

func PointsForChoice(myChoice string) int {
	switch myChoice {
	case "X":
		return 1
	case "Y":
		return 2
	case "Z":
		return 3
	default:
		panic("bad choice: " + myChoice)
	}
}

func GetBattleScore(me string, result BattleResult) int {
	var score int = 0
	switch result {
	case Win:
		score += 6
	case Draw:
		score += 3
	case Lose:
	}
	score += PointsForChoice(me)
	return score
}
func P1GetChoices(line string) (string, string) {
	elf, me := ParseLine(line)
	return elf, me
}

func P2GetChoices(line string) (string, string) {
	// me: x = lose, y = draw, z = win
	var needTo BattleResult
	elf, myGoal := ParseLine(line)
	switch myGoal {
	case "X":
		needTo = Lose
	case "Y":
		needTo = Draw
	case "Z":
		needTo = Win
	}

	for _, choice := range []string{"X", "Y", "Z"} {
		result := DoBattle(elf, choice)
		if result == needTo {
			return elf, choice
		}
	}
	panic("bad")
}

func ParseLine(line string) (string, string) {
	elf := string(line[0])
	me := string(line[len(line)-1])
	return elf, me
}

func main() {
	f, err := os.ReadFile(os.Args[1])
	Fatal(err, "failed to open file")
	data := string(f)
	var lines []string = strings.Split(data, "\n")
	lines = lines[:len(lines)-1] // remove last empty string

	/*
		Part 1
	*/
	var score int = 0
	for _, line := range lines {
		elf, me := P1GetChoices(line)
		result := DoBattle(elf, me)
		score += GetBattleScore(me, result)
	}
	fmt.Println("Part1:", score)

	/*
		Part 2
	*/
	score = 0
	for _, line := range lines {
		elf, me := P2GetChoices(line)
		result := DoBattle(elf, me)
		score += GetBattleScore(me, result)
	}
	fmt.Println("Part2:", score)
}
