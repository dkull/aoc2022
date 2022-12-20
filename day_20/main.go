package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
Utils
*/

type Num interface {
	int | int64
}

func Fatal(err error) {
	if err != nil {
		panic(err)
	}
}

func Modulus[T Num](a, b T) T {
	// return the modulus of a and b
	return ((a % b) + b) % b
}

/*
Structures
*/

type SeqNum[T Num] struct {
	originalOffset T
	value          T
}

type RingBuffer[T Num] struct {
	insertIdxs T
	numbers    []SeqNum[T]
	index      T
}

func (rb *RingBuffer[T]) Step(steps T) {
	/*correct and optimal handling of moving around the ring buffer
	with both positive and negative numbers of any magnitude*/
	rb.index = Modulus(steps+rb.index, T(len(rb.numbers)))
}

func (rb *RingBuffer[T]) Add(num T) {
	newSeqNum := new(SeqNum[T])
	newSeqNum.originalOffset = rb.insertIdxs
	newSeqNum.value = num
	rb.numbers = append(rb.numbers, *newSeqNum)
	rb.insertIdxs++
}

func (rb *RingBuffer[T]) FindInsertIdx(originalOffset T) {
	for rb.numbers[rb.index].originalOffset != originalOffset {
		rb.Step(1)
	}
}

func (rb *RingBuffer[T]) FindValue(value T) {
	for rb.numbers[rb.index].value != value {
		rb.Step(1)
	}
}

func (rb *RingBuffer[T]) ReadValue() T {
	// read the value at the current index
	return rb.numbers[rb.index].value
}

func (rb *RingBuffer[T]) ShuffleValue() {
	// get the value of the element at the current index
	value := rb.numbers[rb.index].value
	item := rb.numbers[rb.index]

	newLoc := Modulus(value+rb.index, T(len(rb.numbers)-1))
	newSeqNums := new([]SeqNum[T])
	if newLoc > rb.index {
		*newSeqNums = append(*newSeqNums, rb.numbers[:rb.index]...)
		*newSeqNums = append(*newSeqNums, rb.numbers[rb.index+1:newLoc+1]...)
		*newSeqNums = append(*newSeqNums, item)
		*newSeqNums = append(*newSeqNums, rb.numbers[newLoc+1:]...)
	} else {
		*newSeqNums = append(*newSeqNums, rb.numbers[:newLoc]...)
		*newSeqNums = append(*newSeqNums, item)
		*newSeqNums = append(*newSeqNums, rb.numbers[newLoc:rb.index]...)
		*newSeqNums = append(*newSeqNums, rb.numbers[rb.index+1:]...)
	}
	rb.numbers = *newSeqNums
}

/*
Functions
*/

func run[T Num](rb RingBuffer[T], mixtimes int) T {
	for i := 0; i < mixtimes; i++ {
		for elemIdx := T(0); elemIdx < rb.insertIdxs; elemIdx++ {
			// find the element in the RingBuffer
			rb.FindInsertIdx(elemIdx)
			rb.ShuffleValue()
		}
	}
	// the answer
	rb.FindValue(0)
	sum := T(0)
	for _, steps := range []T{1000, 1000, 1000} {
		rb.Step(steps)
		sum += rb.ReadValue()
	}
	return sum
}

func main() {
	// read in file using os.ReadFile from Argv[1]
	data, err := os.ReadFile(os.Args[1])
	Fatal(err)
	// split the file into lines
	lines := strings.Split(string(data), "\n")
	// trim the last empty line
	lines = lines[:len(lines)-1]
	// create a new RingBuffer
	rb := RingBuffer[int64]{}
	// loop through the lines
	for _, line := range lines {
		// parse the line as a number
		num, err := strconv.Atoi(line)
		Fatal(err)
		// add the number to the RingBuffer
		rb.Add(int64(num))
	}

	now := time.Now()
	fmt.Println("Part 1:", run(rb, 1))
	fmt.Println("Took", time.Since(now).Round(time.Millisecond))

	for idx := range rb.numbers {
		rb.numbers[idx].value = int64(rb.numbers[idx].value * 811589153)
	}
	now = time.Now()
	fmt.Println("Part 2:", run(rb, 10))
	fmt.Println("Took", time.Since(now).Round(time.Millisecond))
}
