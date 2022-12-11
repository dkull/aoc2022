package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Op struct {
	Type  rune
	Value *int
}

type Item struct {
	Value int
	Ops   []Op
}

/*
do all item ops in modulus
Op value nil means the use the current value ('old' in task)
*/
func (i *Item) ModulusItem(mod int) int {
	result := i.Value
	for _, op := range i.Ops {
		var value int
		if op.Value == nil {
			value = result
		} else {
			value = *op.Value
		}
		switch op.Type {
		case '+':
			result += value
		case '*':
			result *= value
		}
		result = result % mod
	}
	return result
}

type Monkey struct {
	Id              int
	Items           []Item
	Operation       []string
	TestDivisibleBy int
	ThrowToTrue     int
	ThrowToFalse    int
	InspectionCount int
}

/*
monkey goes over all items.
calculate operation outcome for each line.
'old' means the current item, and the possible operations are add and multiply.
it can be 'old + x' or 'old + old'.
then divide that number by 3 and round down.
check the TestDivisibleBy number. if it's divisible, throw to ThrowToTrue, else throw to ThrowToFalse.
*/
func (m *Monkey) MonkeyTurnP1(monkeys *[]Monkey) {
	// for each item
	for _, item := range m.Items {
		(*m).InspectionCount++

		// calculate the operation
		newValue := 0

		var varA = item.Value
		var varB, err = strconv.Atoi(m.Operation[4])
		if err != nil {
			varB = item.Value
		}
		switch m.Operation[3] {
		case "+":
			newValue = varA + varB
		case "*":
			newValue = varA * varB
		}
		// divide by 3 and round down
		newValue = newValue / 3
		// check the TestDivisibleBy number
		item.Value = newValue
		if newValue%m.TestDivisibleBy == 0 {
			(*monkeys)[m.ThrowToTrue].Items = append((*monkeys)[m.ThrowToTrue].Items, item)
		} else {
			(*monkeys)[m.ThrowToFalse].Items = append((*monkeys)[m.ThrowToFalse].Items, item)
		}
	}
	// this monkey has thrown all items, so clear the list
	m.Items = []Item{}
}

/*
we store all operations and do them all every time mod X in ModulusItem
*/
func (m *Monkey) MonkeyTurnP2(monkeys *[]Monkey) {
	// for each item
	for _, item := range m.Items {
		(*m).InspectionCount++

		var value *int
		var parsedNum, err = strconv.Atoi(m.Operation[4])
		if err != nil {
			value = nil
		} else {
			value = &parsedNum
		}

		switch m.Operation[3] {
		case "+":
			item.Ops = append(item.Ops, Op{Type: '+', Value: value})
		case "*":
			item.Ops = append(item.Ops, Op{Type: '*', Value: value})
		}

		isDivisible := item.ModulusItem(m.TestDivisibleBy) == 0
		if isDivisible {
			(*monkeys)[m.ThrowToTrue].Items = append((*monkeys)[m.ThrowToTrue].Items, item)
		} else {
			(*monkeys)[m.ThrowToFalse].Items = append((*monkeys)[m.ThrowToFalse].Items, item)
		}
	}
	// this monkey has thrown all items, so clear the list
	m.Items = []Item{}
}

/*
parse a string of the format:
Monkey 0:

	Starting items: 79, 98
	Operation: new = old * 19
	Test: divisible by 23
	  If true: throw to monkey 2
	  If false: throw to monkey 3

Operation can be:

	old + old
	old + 3
	etc.
*/
func ParseMonkeyString(s string) (Monkey, error) {
	m := Monkey{}
	m.Id = 0
	m.Items = []Item{}
	m.Operation = []string{}
	m.TestDivisibleBy = 0
	m.ThrowToTrue = 0
	m.ThrowToFalse = 0

	// split on newlines
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		// skip if empty  line
		if line == "" {
			continue
		}
		// remove newlines from end of line
		line = strings.TrimSuffix(line, "\n")
		// split on colon
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			return m, errors.New("ParseMonkeyString: invalid line: " + line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		// if key starts with 'Monkey', parse the id using Sscanf
		if strings.HasPrefix(key, "Monkey") {
			_, err := fmt.Sscanf(key, "Monkey %d", &m.Id)
			if err != nil {
				return m, err
			}
			continue
		}
		// print out line if value starts with "If false"
		switch key {
		case "Starting items":
			for _, item := range strings.Split(value, ",") {
				i, _ := strconv.Atoi(strings.TrimSpace(item))
				m.Items = append(m.Items, Item{Value: i})
			}
		case "Operation":
			for _, op := range strings.Split(value, " ") {
				m.Operation = append(m.Operation, op)
			}
		case "Test":
			// divisible by 23
			// using Atoi to get the number
			_, err := fmt.Sscanf(value, "divisible by %d", &m.TestDivisibleBy)
			if err != nil {
				return m, errors.New("ParseMonkeyString: invalid test: " + value)
			}
		case "If true":
			_, err := fmt.Sscanf(value, "throw to monkey %d", &m.ThrowToTrue)
			if err != nil {
				return m, errors.New("ParseMonkeyString: invalid throw to true: " + value)
			}
		case "If false":
			_, err := fmt.Sscanf(value, "throw to monkey %d", &m.ThrowToFalse)
			if err != nil {
				return m, errors.New("ParseMonkeyString: invalid throw to false: " + value)
			}
		default:
			return m, errors.New("ParseMonkeyString: invalid key: " + key)
		}
	}
	return m, nil
}

/*
read a text file and split by \n\n, call ParseMonkey with each block
*/
func ParseMonkeysString(s string) ([]Monkey, error) {
	monkeys := []Monkey{}
	lines := strings.Split(s, "\n\n")
	for _, line := range lines {
		m, err := ParseMonkeyString(line)
		Fatal(err)
		monkeys = append(monkeys, m)
	}
	return monkeys, nil
}

func main() {
	// read in a file from path in os.Args[1]
	file, err := ioutil.ReadFile(os.Args[1])
	Fatal(err)
	// parse the file
	// == Part 1 ==
	monkeys, err := ParseMonkeysString(string(file))
	Fatal(err)
	// run the simulation for 20 rounds
	for i := 0; i < 20; i++ {
		for i, monkey := range monkeys {
			monkey.MonkeyTurnP1(&monkeys)
			monkeys[i] = monkey
		}
	}
	// sort monkeys by highest InspectionCount
	// multiply top 2 monkey inspection counts and Print as Part1
	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].InspectionCount > monkeys[j].InspectionCount
	})
	fmt.Printf("Part1: %d\n", monkeys[0].InspectionCount*monkeys[1].InspectionCount)

	// == Part 2 ==
	monkeys, err = ParseMonkeysString(string(file))
	Fatal(err)
	// run the simulation for 10000 rounds
	for i := 0; i < 10000; i++ {
		for i, monkey := range monkeys {
			monkey.MonkeyTurnP2(&monkeys)
			monkeys[i] = monkey
		}
	}
	// sort monkeys by highest InspectionCount
	// multiply top 2 monkey inspection counts and Print as Part2
	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].InspectionCount > monkeys[j].InspectionCount
	})
	fmt.Printf("Part2: %d\n", monkeys[0].InspectionCount*monkeys[1].InspectionCount)
}
