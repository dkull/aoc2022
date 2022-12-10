package main

import (
	"bufio"
	"fmt"
	"os"
)

func Fail(err error) {
	if err != nil {
		panic(err)
	}
}

func ReadFileToLines(filename string) []string {
	file, err := os.Open(filename)
	Fail(err)
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	Fail(scanner.Err())
	return lines
}

/*
this function takes a string array of lines and returns the summed signal
strengths.
parse the instructions:

	"addx <value>" // value can be negative or 0. takes 2 cycles to run.
	"noop" // takes 1 cycle to run.

have a register X and a totalCycles counter.
loop over all lines. inside the loop parse the instruction.
if we have an addx instruction, store the 'cycles' value of 2 and the Value.
if we have an noop instruction, store the 'cycles' value of 1 and the Value 0.

inside the 'cycles' loop:
every cycle run a check for (totalCycles + 20) % 40 == 0. if that modulus operation is true, add (totalCycles * X) to get a signalStrength and add it to sumSignalStength.
return sumSignalStrength.
*/
func RunLines(lines []string) int {
	sumSignalStrength := 0
	X := 1
	totalCycles := 0
	for _, line := range lines {
		cycles := 1
		value := 0
		if line[0:4] == "addx" {
			cycles = 2
			fmt.Sscanf(line, "addx %d", &value)
		}
		for i := 0; i < cycles; i++ {
			DrawPixel(totalCycles, X) // Part2
			totalCycles++
			if (totalCycles+20)%40 == 0 {
				sumSignalStrength += totalCycles * X
			}
		}
		X += value
	}
	return sumSignalStrength
}

func DrawPixel(cycle int, X int) {
	const cols = 40
	var pixelCol int = cycle % cols
	var spriteCol int = X % cols
	// if spriteCol is -1..1 away from pixelCol then print '#' else '.'
	if pixelCol == 0 {
		fmt.Println()
	}
	if spriteCol >= pixelCol-1 && spriteCol <= pixelCol+1 {
		fmt.Print("#")
	} else {
		fmt.Print(".")
	}
}

/*
parse lines of file Argv[1] and print the result as Part1:
*/
func main() {
	lines := ReadFileToLines(os.Args[1])
	result := RunLines(lines)
	println("Part1:", result)
	// Part2 is printed first, it's inlined
}
