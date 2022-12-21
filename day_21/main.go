package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/JohnCGriffin/overflow"
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

type Result struct {
	Left    int
	Right   int
	Failure bool // mark a result that we can't trust. eg. float division happened
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
Functions
*/

/*
given a hashmap of map[string]Monkey, resolve all all the monkeys Expression
values recusively until all monkeys have a single value
*/
func resolveMonkeys1(monkeys map[string]*Monkey, target string, results *map[string]int) {
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
		(*results)[targetMonkey.Name] = int(number)
	} else {
		// if the monkey has more than one value, it is an expression
		// and we need to resolve it
		left := targetMonkey.Expression[0]
		operator := targetMonkey.Expression[1]
		right := targetMonkey.Expression[2]
		// if the left or right values are not numbers, they are monkeys
		// and we need to resolve them first
		resolveMonkeys1(monkeys, left, results)
		resolveMonkeys1(monkeys, right, results)
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

func resolveMonkeys2(monkeys map[string]*Monkey, target string, results *map[string]int, result *Result) {
	targetMonkey := monkeys[target]
	if result.Failure {
		return
	}

	// skip monkeys that have already been resolved
	if _, ok := (*results)[targetMonkey.Name]; ok {
		return
	}

	if len(targetMonkey.Expression) == 1 {
		// if the monkey has a single value, it is a number
		//fmt.Println("resolving monkeys target:", targetMonkey.Name, "value:", targetMonkey.Expression[0])
		number, err := strconv.Atoi(targetMonkey.Expression[0])
		Fatal(err)
		// and we can add it to the results map
		(*results)[targetMonkey.Name] = int(number)
	} else {
		// if the monkey has more than one value, it is an expression
		// and we need to resolve it
		left := targetMonkey.Expression[0]
		operator := targetMonkey.Expression[1]
		right := targetMonkey.Expression[2]
		// if the left or right values are not numbers, they are monkeys
		// and we need to resolve them first
		resolveMonkeys2(monkeys, left, results, result)
		resolveMonkeys2(monkeys, right, results, result)
		if result.Failure {
			return
		}
		leftValue := (*results)[left]
		rightValue := (*results)[right]
		switch operator {
		case "=":
			result.Left = leftValue
			result.Right = rightValue
		case "+":
			sum, ok := overflow.Add(int(leftValue), int(rightValue))
			result.Failure = !ok
			(*results)[targetMonkey.Name] = sum
		case "-":
			sub, ok := overflow.Sub(int(leftValue), int(rightValue))
			result.Failure = !ok
			(*results)[targetMonkey.Name] = sub
		case "*":
			prod, ok := overflow.Mul(int(leftValue), int(rightValue))
			result.Failure = !ok
			(*results)[targetMonkey.Name] = prod
		case "/":
			// we can't do this operation if the two numbers do not divide cleanly
			if rightValue != 0 && leftValue%rightValue == 0 {
				(*results)[targetMonkey.Name] = leftValue / rightValue
			} else {
				result.Failure = true
			}
		}
	}
}

/*
return true if number a is between one side of <arg> and b is on another
example a = 5, b = 10, arg = 7 returns true
*/
func between(a, b, arg float64) bool {
	return (a < arg && arg < b) || (b < arg && arg < a)
}

func GradientLocalMinimaSeeker(monkeys map[string]*Monkey, results *map[string]int, result *Result) {
	min := int(math.MinInt32)
	stepSize := int(1000000)
	prevRatio := float64(0)
	value := min
	backtracked := false

	for {
		// resolve the monkeys
		results := make(map[string]int)
		results["humn"] = value
		result := Result{}
		resolveMonkeys2(monkeys, "root", &results, &result)

		// if we failed at math, try a tiny step forward until we can figure out a gradient again
		if result.Failure {
			value += 1
			continue
		}

		left := result.Left
		right := result.Right
		fmt.Println("======", value, "left:", left, "right:", right, "======")

		// check if we won
		if left == right && result.Failure == false {
			return
		}

		// compute the gradient
		ratio := (float64(left) - float64(right)) / float64(right)
		if prevRatio == 0.0 {
			prevRatio = ratio
			continue
		}

		// check if we jumped over
		flipped := between(ratio, prevRatio, float64(0.0))
		fmt.Println("prevRatio", prevRatio, "ratio", ratio, "backtracked", backtracked, "flipped", flipped)

		if flipped {
			// if the ratio flipped, we need to start pinpointing
			// go back to previous value for this, and use a smaller step size
			value = value - stepSize
			stepSize = stepSize / 2
			// try this step again, but with a smaller step
			fmt.Println("flipped, backtracking")
			backtracked = true
			continue
		}
		backtracked = false

		// finalize this step
		stepSize *= int(math.Ceil(math.Abs(ratio)))
		prevRatio = ratio
		value += stepSize

		fmt.Println("newStepSize", stepSize)
	}
}

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
	resolveMonkeys1(monkeys, "root", &results)
	// print the result of "root" monkey as Part1:
	println("Part1:", results["root"])

	/* Part 2 */
	result := Result{}
	monkeys["root"].Expression[1] = "="
	GradientLocalMinimaSeeker(monkeys, &results, &result)
}
