package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

func AbsDelta(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}

func Pow(b, p int) (result int) {
	result = 1
	for i := 0; i < p; i++ {
		result *= b
	}
	return result
}

func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

/*
decode format: 2=0=
the rightmost place is the 5s place,
one left of that is the 5*5s place,
etc.
= is -2
- is -1
0, 1, and 2 are themselves
*/
func Decode(inp string) int {
	sum := 0
	for i := 0; i < len(inp); i++ {
		reverseIdx := len(inp) - 1 - i
		char := inp[reverseIdx]
		multiplier := Pow(5, i)
		switch char {
		case '=':
			sum += -2 * multiplier
		case '-':
			sum += -1 * multiplier
		case '0':
			sum += 0 * multiplier
		case '1':
			sum += 1 * multiplier
		case '2':
			sum += 2 * multiplier
		}
	}
	return sum
}

func Encode(inp int) (out string) {
	symbols := make(map[string]int)
	symbols["="] = -2
	symbols["-"] = -1
	symbols["0"] = 0
	symbols["1"] = 1
	symbols["2"] = 2

	topPow := 0
	absDelta := math.MaxInt
	for pow := 30; pow >= 0; pow-- {
		bestVal := 0
		bestSym := "N/A"
		for sym, val := range symbols {
			topPow = Pow(5, pow)
			if AbsDelta(topPow*val, inp) <= absDelta {
				absDelta = AbsDelta(topPow*val, inp)
				bestVal = val * topPow
				bestSym = sym
			}
		}
		if bestSym == "N/A" {
			continue
		}
		out += bestSym
		inp -= bestVal
	}

	// trim leading '0' in out
	for out[0] == '0' {
		out = out[1:]
	}

	return
}

func main() {
	// read data using os.ReadFile from Argv[1]
	data, err := os.ReadFile(os.Args[1])
	Fatal(err)
	// test decoder:
	lines := strings.Split(string(data), "\n")
	p1Sum := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		p1Sum += Decode(line)
	}
	fmt.Println("Part1Raw:", p1Sum)
	fmt.Println("Part1Encoded:", Encode(p1Sum))
}
