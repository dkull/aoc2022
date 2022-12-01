package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func Fatal(err error, msg string) {
	if err != nil {
		panic(msg)
	}
}

type Backpack struct {
	rations []int
}

func (b *Backpack) FromLines(lines []string) []string {
	for i, line := range lines {
		if line == "" {
			return lines[i+1:]
		}
		calories, err := strconv.Atoi(line)
		Fatal(err, "failed to parse integer")
		b.rations = append(b.rations, calories)
	}
	return []string{}
}

func (b Backpack) CalorieSum() int {
	var sum int = 0
	for _, r := range b.rations {
		sum = sum + r
	}
	return sum
}

func PopulateBackpacks(lines []string) []Backpack {
	var backpacks []Backpack
	var backpack *Backpack = new(Backpack)

	for {
		lines = backpack.FromLines(lines)
		backpacks = append(backpacks, *backpack)
		backpack = new(Backpack)
		if len(lines) == 0 {
			break
		}
	}
	return backpacks
}

func main() {
	filePath := os.Args[1]

	data, err := os.ReadFile(filePath)
	Fatal(err, "failed to find input file")

	lines := strings.Split(string(data), "\n")

	backpacks := PopulateBackpacks(lines)

	// Part 1
	var highest int = 0
	for _, bp := range backpacks {
		if calories := bp.CalorieSum(); calories > highest {
			highest = calories
		}
	}
	fmt.Println("Part1:", highest)

	// Part 2
	var calories []int
	for _, bp := range backpacks {
		calories = append(calories, bp.CalorieSum())
	}
	sort.Sort(sort.Reverse(sort.IntSlice(calories)))
	var top3Sum int = 0
	for _, cal := range calories[0:3] {
		top3Sum = top3Sum + cal
	}
	fmt.Println("Part2:", top3Sum)
}
