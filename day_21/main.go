package main

import (
	"os"
	"strconv"
	"strings"
)

/*
Utils
*/

func Fatal(err error) {
	if err != nil {
		panic(err)
	}
}

/*
Datastructures
*/

type Monkey struct {
	Name       string
	Expression []string
}

/*
Parsing
*/

/*
Parse lines into Monkey structs:
root: pppw + sjmn
cczh: sllz * lgvd
dbpl: 5
cczh: sllz / lgvd
ptdq: humn - dvpt
*/
func parseMonkey(line string) Monkey {
	var monkey Monkey
	tokens := strings.Split(line, " ")
	monkey.Name = tokens[0][:len(tokens[0])-1]
	monkey.Expression = tokens[1:]
	return monkey
}

/*
given a hashmap of map[string]Monkey, resolve all all the monkeys Expression
values recusively until all monkeys have a single value
*/
func resolveMonkeys(monkeys map[string]*Monkey, target string, results *map[string]int) {
	targetMonkey := monkeys[target]

	// skip monkeys that have already been resolved
	if _, ok := (*results)[targetMonkey.Name]; ok {
		return
	}

	if len(targetMonkey.Expression) == 1 {
		// if the monkey has a single value, it is a number
		number, err := strconv.Atoi(targetMonkey.Expression[0])
		Fatal(err)
		// and we can add it to the results map
		(*results)[targetMonkey.Name] = number
	} else {
		// if the monkey has more than one value, it is an expression
		// and we need to resolve it
		left := targetMonkey.Expression[0]
		operator := targetMonkey.Expression[1]
		right := targetMonkey.Expression[2]
		// if the left or right values are not numbers, they are monkeys
		// and we need to resolve them first
		resolveMonkeys(monkeys, left, results)
		resolveMonkeys(monkeys, right, results)
		switch operator {
		case "+":
			(*results)[targetMonkey.Name] = (*results)[left] + (*results)[right]
		case "-":
			(*results)[targetMonkey.Name] = (*results)[left] - (*results)[right]
		case "*":
			(*results)[targetMonkey.Name] = (*results)[left] * (*results)[right]
		case "/":
			(*results)[targetMonkey.Name] = (*results)[left] / (*results)[right]
		}
	}
}

/*
Functions
*/

func main() {
	// read file from Argv[1] using os.ReadFile
	fileData, err := os.ReadFile(os.Args[1])
	Fatal(err)
	// split it into lines
	lines := strings.Split(string(fileData), "\n")
	// remove last empty line
	lines = lines[:len(lines)-1]
	// parse lines into Monkey structs
	monkeys := make(map[string]*Monkey)
	for _, line := range lines {
		monkey := parseMonkey(line)
		monkeys[monkey.Name] = &monkey
	}
	// resolve all the monkeys
	results := make(map[string]int)
	resolveMonkeys(monkeys, "root", &results)
	// print the result of "root" monkey as Part1:
	println("Part1:", results["root"])
}
