package main

import (
	"fmt"
	"os"
	"strings"
)

/*
This function reads a line and splits it in half.
It then finds a common character on both sides of the line.
There is only 1 common character per line.
Each common lowercase character gets a score of 1-26 for a-z.
Each uppercase character gets a score of 27-52 for A-Z.
*/
func LineScorer(line string) int {
	// Split the line in half
	half := len(line) / 2
	left := line[:half]
	right := line[half:]

	// Find the common character
	var common rune
	for _, l := range left {
		for _, r := range right {
			if l == r {
				common = l
				break
			}
		}
	}

	// Calculate the score
	var score int
	switch {
	case common >= 'a' && common <= 'z':
		score = int(common-'a') + 1
	case common >= 'A' && common <= 'Z':
		score = int(common-'A') + 27
	}

	return score
}

func Fatal(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}

/*
count characters in each item separately.
then count how many maps contain the character.
return characters that are present in all maps
do not use indexing
*/
func Intersection(items []string) string {
	// Count characters in each item
	var maps []map[rune]int
	for _, item := range items {
		m := make(map[rune]int)
		for _, c := range item {
			m[c]++
		}
		maps = append(maps, m)
	}

	// Count how many maps contain the character
	var count map[rune]int = make(map[rune]int)
	for _, m := range maps {
		for c := range m {
			if _, ok := count[c]; ok {
				count[c]++
			} else {
				count[c] = 1
			}
		}
	}

	// Return characters that are present in all maps
	var result string
	for c, n := range count {
		if n == len(maps) {
			result += string(c)
		}
	}
	return result
}

/*
Read lines into groups of 3 lines.
Find the intersection of the 3 lines.
The intersection is the group key.

Example input:
vJrwpWtwJgWrhcsFMMfFFhFp
jqHRNqRjqzjGDLGLrsFMfFZSrLrFZsSL
wMqvLMZHhHMvwLHjbvcjnnSBnvTQFn
ttgJtRGJQctTZtZT
PmmdzqPrVvPwwTWBwg
CrZsJsPPZsGzwwsLwLmpwMDw
Output: []{'r', 'Z'}
*/
func LineGrouper(lines []string) []rune {
	var groups []rune
	for i := 0; i < len(lines); i += 3 {
		groups = append(groups, []rune(Intersection(lines[i:i+3]))...)
	}
	return groups
}

/*
Each line is fed into LineGrouper to find groups of 3 lines.
Each group key is then scored a-z=1-26 and A-Z=27-25.
The scores are calculated without LineScorer, and on the group key.
*/
func GroupScorer(lines []string) int {
	var total int
	groups := LineGrouper(lines)
	for _, group := range groups {
		switch {
		case group >= 'a' && group <= 'z':
			total += int(group-'a') + 1
		case group >= 'A' && group <= 'Z':
			total += int(group-'A') + 27
		}
	}
	return total
}

/*
This function reads a file, then calls LineScorer with each line and adds the score to the total.
Return the score.
The code is compact.
*/
func FileScorer(filename string) int {
	var total int
	data, err := os.ReadFile(filename)
	Fatal(err, "Failed to open file")
	var lines []string = strings.Split(string(data), "\n")
	for _, line := range lines {
		total += LineScorer(line)
	}
	return total
}

/*
The function ReadFileLines reads a file and returns an array of lines.
Remove empty lines.
Manual!
*/
func ReadFileLines(filename string) []string {
	data, err := os.ReadFile(filename)
	Fatal(err, "Failed to open file")
	var lines []string = strings.Split(string(data), "\n")
	var result []string
	for _, line := range lines {
		if line != "" {
			result = append(result, line)
		}
	}
	return result
}

func main() {
	fmt.Println("Part1:", FileScorer(os.Args[1]))
	fmt.Println("Part2:", GroupScorer(ReadFileLines(os.Args[1])))
}
