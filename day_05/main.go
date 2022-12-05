package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Fatal(err error) {
	if err != nil {
		panic(err)
	}
}

func PrettyPrintState(state []string) {
	for _, line := range state {
		fmt.Println(line)
	}
}

/*
Transpose a multiline input, with all lines having the same length.
Only check the length of the first line.
*/
func Transpose(input []string) []string {
	// Create a slice of strings to hold the transposed lines.
	transposed := make([]string, len(input[0]))
	// Loop over the lines.
	for _, line := range input {
		// Loop over the characters in the line.
		for i, char := range line {
			// Add the character to the transposed line.
			transposed[i] += string(char)
		}
	}
	return transposed
}

/*
Split an input file into two parts.
The first part is before an empty line.
The second part is after the empty line.
Remove last line from second part.
*/
func SplitInput(input string) ([]string, []string) {
	// Split the input into lines.
	lines := strings.Split(input, "\n")
	// Find the index of the empty line.
	emptyLineIndex := -1
	for i, line := range lines {
		if line == "" {
			emptyLineIndex = i
			break
		}
	}
	// Split the lines into two parts.
	firstPart := lines[:emptyLineIndex]
	secondPart := lines[emptyLineIndex+1 : len(lines)-1]
	return firstPart, secondPart
}

/*
Remove lines which do not contain a number anywhere.
Reverse each line with a for loop.
Remove the first character.
Remove trailing spaces.
*/
func CleanFirstPart(input []string) []string {
	// Create a slice of strings to hold the cleaned lines.
	cleaned := make([]string, 0)
	// Loop over the lines.
	for _, line := range input {
		// Check if the line contains a number.
		if strings.ContainsAny(line, "0123456789") {
			// Create a slice of runes to hold the reversed line.
			reversed := make([]rune, len(line))
			// Loop over the characters in the line.
			for i, char := range line {
				// Add the character to the reversed line.
				reversed[len(line)-i-1] = char
			}
			// Convert the reversed line to a string.
			reversedString := string(reversed)
			// Remove the first character.
			reversedString = reversedString[1:]
			// Remove trailing spaces.
			reversedString = strings.TrimRight(reversedString, " ")
			// Add the reversed line to the cleaned lines.
			cleaned = append(cleaned, reversedString)
		}
	}
	return cleaned
}

/*
Take instructions from the second part of the input file.
The instructions are of the form:
move X from Y to Z
move 3 from 1 to 3
The first number means how many, from means from which row and to means to which row.
multiContainer is a boolean which indicates if the instructions are for the multi container function
If multiContainer is true call MoveMultipleItemsFromTo, else call MoveSingleItemsFromTo. State is the last argument.
Does not modify its input state.
*/
func MoveItems(instructions []string, state []string, multiContainer bool) []string {
	// Create a slice of strings to hold the new state.
	newState := make([]string, len(state))
	// Copy the state to the new state.
	copy(newState, state)
	// Loop over the instructions.
	for _, instruction := range instructions {
		// Split the instruction into words.
		words := strings.Split(instruction, " ")
		// Get the number of items to move.
		number, err := strconv.Atoi(words[1])
		Fatal(err)
		// Get the row to move from.
		from, err := strconv.Atoi(words[3])
		Fatal(err)
		// Get the row to move to.
		to, err := strconv.Atoi(words[5])
		Fatal(err)
		// Move the items.
		if multiContainer {
			newState = MoveMultipleItemsFromTo(number, from, to, newState)
		} else {
			newState = MoveSingleItemsFromTo(number, from, to, newState)
		}
	}
	return newState
}

/*
Move single items from one row to another.
HowMany means one by one, a repeated action.
From is from which row. To is to which row.
The from and to indexes are 1-based.
*/
func MoveSingleItemsFromTo(howMany, from, to int, state []string) []string {
	//fmt.Println("=== Move", howMany, "from", from, "to", to)
	// Loop over the number of items to move.
	for i := 0; i < howMany; i++ {
		// Get the last item from the from row.
		item := state[from-1][len(state[from-1])-1:]
		// Remove the last item from the from row.
		state[from-1] = state[from-1][:len(state[from-1])-1]
		// Add the item to the to row.
		state[to-1] += item
		// prettyprint the state
		//PrettyPrintState(state)
	}
	return state
}

/*
just from signature
*/
func MoveMultipleItemsFromTo(howMany, from, to int, state []string) []string {
	//fmt.Println("=== Move", howMany, "from", from, "to", to)
	// Get the items to move.
	items := state[from-1][len(state[from-1])-howMany:]
	// Remove the items from the from row.
	state[from-1] = state[from-1][:len(state[from-1])-howMany]
	// Add the items to the to row.
	state[to-1] += items
	// prettyprint the state
	//PrettyPrintState(state)
	return state
}

/*
Return a string concatenated from the last element of each row.
*/
func GetResult(state []string) string {
	// Create a slice of strings to hold the result.
	result := make([]string, 0)
	// Loop over the rows.
	for _, row := range state {
		// Get the last character from the row.
		char := row[len(row)-1:]
		// Add the character to the result.
		result = append(result, char)
	}
	// Join the result into a string.
	return strings.Join(result, "")
}

/*
read input with with os.ReadFile from os.Argv[1]
SplitInput() on input.
Transpose() the first part.
CleanFirstPart() the transposed lines.
First call MoveItems(cleaned, single) and print out the result as Part1.
Then call MoveItems(cleaned, multi) and print out the result as Part2.
Use the same instance of cleaned data on both calls.
*/
func main() {
	// Read the input file.
	input, err := os.ReadFile(os.Args[1])
	Fatal(err)
	// Split the input into two parts.
	firstPart, secondPart := SplitInput(string(input))
	// Transpose the first part.
	firstPartTransposed := Transpose(firstPart)
	// Clean the first part.
	firstPartCleaned := CleanFirstPart(firstPartTransposed)
	// Move the items and print out the result.
	fmt.Println("Part1:", GetResult(MoveItems(secondPart, firstPartCleaned, false)))
	fmt.Println("Part2:", GetResult(MoveItems(secondPart, firstPartCleaned, true)))
}
