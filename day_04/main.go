package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// from signature
func Fatal(e error) {
	if e != nil {
		panic(e)
	}
}

/*
A struct containing two pairs of number ranges. Call them A and B.
*/
type Pairs struct {
	A, B [2]int
}

// Method on Pairs to check if either A contains B or B contains A
// Takes no arguments
func (p Pairs) Contains() bool {
	return p.A[0] <= p.B[0] && p.A[1] >= p.B[1] || p.B[0] <= p.A[0] && p.B[1] >= p.A[1]
}

// Method on pairs to check if A overlaps B in any place
func (p Pairs) Overlaps() bool {
	return p.A[0] <= p.B[0] && p.A[1] >= p.B[0] || p.B[0] <= p.A[0] && p.B[1] >= p.A[0]
}

// PairsParser returns a Pairs struct containing the two pairs of numbers.
func PairsParser(line string) Pairs {
	var p Pairs
	parts := strings.Split(line, ",")
	for i, part := range parts {
		nums := strings.Split(part, "-")
		for j, num := range nums {
			n, _ := strconv.Atoi(num)
			if i == 0 {
				p.A[j] = n
			} else {
				p.B[j] = n
			}
		}
	}
	return p
}

/*
Part 1 code.
Just call Contains() on each line.
Sum up the number of trues.
*/
func Part1(pairs []Pairs) int {
	count := 0
	for _, p := range pairs {
		if p.Contains() {
			count++
		}
	}
	return count
}

/*
Part 2 code
Just call Overlaps() on each line.
Sum up the number of trues.
*/
func Part2(pairs []Pairs) int {
	count := 0
	for _, pair := range pairs {
		if pair.Overlaps() {
			count++
		}
	}
	return count
}

/*
Example file:
2-4,6-8
2-3,4-5
5-7,7-9
*/
func main() {
	// read a file from os.Argv[1]. using os.ReadFile
	data, e := os.ReadFile(os.Args[1])
	Fatal(e)
	lines := strings.Split(string(data), "\n")
	// remove the last line if it's empty
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	// call PairsParser on each line
	pairs := make([]Pairs, len(lines))
	for i, line := range lines {
		pairs[i] = PairsParser(line)
	}

	fmt.Println(Part1(pairs))
	fmt.Println(Part2(pairs))
}
