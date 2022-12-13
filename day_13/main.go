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

func Max[T int](a, b T) T {
	if a > b {
		return a
	}
	return b
}

type Result int

const (
	Unknown Result = iota
	Right
	Wrong
)

/*
turn a string "[1,[2,3,4,[5,6,7]],8,9]" into a list of {}interface{}s.
eventually it will be a list of integers and lists of integers, ect.
convert the numbers to integers, and the lists to []interface{}
parse it recursively. don't use ParseElement
Largely handwritten
*/
func Parse(list string) (int, []interface{}) {
	var elements []interface{}
	var activeElem string
	for i := 1; i < len(list); i++ {
		switch list[i] {
		case '[':
			var newItem []interface{}
			var iPlus int
			iPlus, newItem = Parse(list[i:])
			i += iPlus
			elements = append(elements, newItem)
		case ']':
			if activeElem != "" {
				num, err := strconv.Atoi(activeElem)
				Fatal(err)
				elements = append(elements, num)
			}
			return i, elements
		case ',':
			if activeElem != "" {
				num, err := strconv.Atoi(activeElem)
				Fatal(err)
				elements = append(elements, num)
			}
			activeElem = ""
		default:
			activeElem += string(list[i])
		}
	}
	return len(list), elements
}

/*
Handwritten
*/
func Extract(list string) []string {
	var result []string
	var current []rune
	var depth int
	for _, c := range list {
		if c == '[' {
			depth++
		} else if c == ']' {
			depth--
		}
		if depth == 1 && c == ',' {
			result = append(result, string(current))
			current = []rune{}
		} else {
			current = append(current, c)
		}
	}
	if len(current) > 0 {
		result = append(result, string(current))
	}
	return result
}

/*
given two lists of lists of strings representing arrays.
recurse element by element and comparing them.
Handwritten
*/
func Compare(a, b []interface{}) Result {
	for i := 0; i < Max(len(a), len(b)); i++ {
		if len(a) == i {
			return Right
		}
		if len(b) == i {
			return Wrong
		}
		aa := a[i]
		bb := b[i]
		switch aa.(type) {
		case int:
			switch bb.(type) {
			case int:
				if aa.(int) > bb.(int) {
					return Wrong
				} else if aa.(int) < bb.(int) {
					return Right
				}
			case []interface{}:
				newA := []interface{}{aa}
				result := Compare(newA, bb.([]interface{}))
				if result == Wrong || result == Right {
					return result
				}
			}
		case []interface{}:
			switch bb.(type) {
			case int:
				newB := []interface{}{bb}
				result := Compare(aa.([]interface{}), newB)
				if result == Wrong || result == Right {
					return result
				}
			case []interface{}:
				result := Compare(aa.([]interface{}), bb.([]interface{}))
				if result == Wrong || result == Right {
					return result
				}
			}
		}
	}
	// 5611 is too high
	// 4753 is too low
	return Unknown
}

/*
load file using os.ReadFile from Argv
run Extract on the first line
and print each output line from Extract
Handwritten
*/
func main() {
	fname := os.Args[1]
	data, err := os.ReadFile(fname)
	Fatal(err)
	scoreP1 := 0
	scoreP2_2 := 1 // start at 1 because ths is the first index
	scoreP2_6 := 2 // start at 1+1 because [[2]] is smaller anyway

	_, line2 := Parse("[[2]]") //part2
	_, line6 := Parse("[[6]]") //part2

	// split the file by two empty lines
	// each group will contain two lines
	// feed both lines to Extract, and then to Compare
	for i, group := range strings.Split(string(data), "\n\n") {
		// part1
		lines := strings.Split(group, "\n")[0:2]
		_, a := Parse(lines[0])
		_, b := Parse(lines[1])
		res := Compare(a, b)
		if res == Right {
			scoreP1 += i + 1
		}

		//part2
		if compare := Compare(a, line2); compare == Right {
			scoreP2_2 += 1
		}
		if compare := Compare(b, line2); compare == Right {
			scoreP2_2 += 1
		}
		if compare := Compare(a, line6); compare == Right {
			scoreP2_6 += 1
		}
		if compare := Compare(b, line6); compare == Right {
			scoreP2_6 += 1
		}

	}
	fmt.Println("Part1:", scoreP1)
	fmt.Println("Part2:", scoreP2_2*scoreP2_6)
}
