package main

import (
	"fmt"
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

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Modulus(a, b int) int {
	// return the modulus of a and b
	modSteps := ((a % b) + b) % b
	return modSteps
	//return Abs((a + modSteps) % b)
}

/*
Structures
*/

type SeqNum struct {
	originalOffset int
	value          int
}

type RingBuffer struct {
	insertIdxs int
	numbers    []SeqNum
	index      int
}

func (rb *RingBuffer) Step(steps int) {
	/*correct and optimal handling of moving around the ring buffer
	with both positive and negative numbers of any magnitude*/
	//modSteps := (steps%len(rb.numbers) + steps) % len(rb.numbers)
	//rb.index = (rb.index + modSteps) % len(rb.numbers)
	rb.index = Modulus(steps+rb.index, len(rb.numbers))
}

func (rb *RingBuffer) Add(num int) {
	rb.numbers = append(rb.numbers, SeqNum{
		originalOffset: rb.insertIdxs,
		value:          num,
	})
	rb.insertIdxs++
}

func (rb *RingBuffer) FindInsertIdx(originalOffset int) {
	for rb.numbers[rb.index].originalOffset != originalOffset {
		rb.Step(1)
	}
}

func (rb *RingBuffer) FindValue(value int) {
	for rb.numbers[rb.index].value != value {
		rb.Step(1)
	}
}

func (rb *RingBuffer) ReadValue() int {
	// read the value at the current index
	return rb.numbers[rb.index].value
}

func (rb *RingBuffer) ShuffleValue() {
	// get the value of the element at the current index
	value := rb.numbers[rb.index].value
	item := rb.numbers[rb.index]
	//newLoc := Modulus(value, len(rb.numbers), &rb.index)

	newLoc := Modulus(value+rb.index, len(rb.numbers)-1)
	fmt.Println("> Shuffling", value, "to", newLoc, "from", rb.index, "while at", rb.index)
	//rb.Step(value)
	if newLoc > rb.index {
		newBuffer2 := [][]SeqNum{
			rb.numbers[:rb.index], rb.numbers[rb.index+1 : newLoc+1], {item}, rb.numbers[newLoc+1:],
		}
		//fmt.Println("newBuffer2GT:", newBuffer2)
		// concatenate the slices
		resultBuffer := make([]SeqNum, 0, len(rb.numbers))
		for _, slice := range newBuffer2 {
			resultBuffer = append(resultBuffer, slice...)
		}
		rb.numbers = resultBuffer
	} else {
		//fmt.Println("item", item, rb.index, newLoc)
		newBuffer2 := [][]SeqNum{
			rb.numbers[:newLoc], {item}, rb.numbers[newLoc:rb.index], rb.numbers[rb.index+1:],
		}
		//fmt.Println("newBuffer2LT:", newBuffer2)
		// concatenate the slices
		resultBuffer := make([]SeqNum, 0, len(rb.numbers))
		for _, slice := range newBuffer2 {
			resultBuffer = append(resultBuffer, slice...)
		}
		rb.numbers = resultBuffer
	}
}

/*
Functions
*/

func part1(rb RingBuffer) int {
	for elemIdx := 0; elemIdx < rb.insertIdxs; elemIdx++ {
		// find the element in the RingBuffer
		rb.FindInsertIdx(elemIdx)
		/*fmt.Println("moving:", rb.numbers[rb.index].value)
		fmt.Print("Before: ")
		for _, elem := range rb.numbers {
			fmt.Print(elem.value, " ")
		}
		fmt.Println()*/
		rb.ShuffleValue()
		/*fmt.Print(" After: ")
		for _, elem := range rb.numbers {
			fmt.Print(elem.value, " ")
		}
		fmt.Println()*/
	}
	// the answer
	rb.FindValue(0)
	sum := 0
	for _, steps := range []int{1000, 1000, 1000} {
		rb.Step(steps)
		sum += rb.ReadValue()
	}
	return sum
}

/*
21725 is too high
*/
func main() {
	// read in file using os.ReadFile from Argv[1]
	data, err := os.ReadFile(os.Args[1])
	Fatal(err)
	// split the file into lines
	lines := strings.Split(string(data), "\n")
	// trim the last empty line
	lines = lines[:len(lines)-1]
	// create a new RingBuffer
	rb := RingBuffer{}
	// loop through the lines
	for _, line := range lines {
		// parse the line as a number
		num, err := strconv.Atoi(line)
		Fatal(err)
		// add the number to the RingBuffer
		rb.Add(num)
	}
	fmt.Println("Part 1:", part1(rb))
}
