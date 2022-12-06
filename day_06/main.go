package main

import (
	"bufio"
	"fmt"
	"os"
)

// from signature
func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

// from signature
func ReadFileToLines(filename string) []string {
	file, err := os.Open(filename)
	Panic(err)
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	Panic(scanner.Err())

	return lines
}

/*
Iterate over characters in a line, with a sliding window of length <seqLen>.
Use a single <seqLen> sized buffer to store the current sequence. Use modulus to assign into it.
Collect all 4 characters in the window in a map.
For each window, create map to check if all the characters inside it are different from each other.
If the map has a length of <seqLen>, then all characters are unique.
If they are, return the index of the last character in the window.
*/
func FindFirstUniqueChar(line string, seqLen int) int {
	// Initialize the buffer
	buffer := make([]byte, seqLen)
	for i := 0; i < seqLen; i++ {
		buffer[i] = line[i]
	}

	// Iterate over the line
	for i := seqLen; i < len(line); i++ {
		// Check if all characters in the buffer are unique
		charMap := make(map[byte]bool)
		for _, char := range buffer {
			charMap[char] = true
		}
		if len(charMap) == seqLen {
			return i
		}

		// Add the next character to the buffer
		buffer[i%seqLen] = line[i]
	}
	return -1
}

/*
Read file using ReadFileToLines. Path is in os.Args[1].
For each line, call FindFirstUniqueChar and print the line with the result with prefix Part1.
Use seqLen 4 for Part1.
Then use seqLen 14 for Part2.
*/
func main() {
	lines := ReadFileToLines(os.Args[1])
	for _, line := range lines {
		fmt.Printf("Part1: %s %d\n", line, FindFirstUniqueChar(line, 4))
		fmt.Printf("Part2: %s %d\n", line, FindFirstUniqueChar(line, 14))
	}
}
