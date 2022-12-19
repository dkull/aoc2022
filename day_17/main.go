package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func contains[T string](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func SlicesEqual[T comparable](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

type Pair[T any] struct {
	a, b T
}

// Shape

type Shape struct {
	b [][]byte
}

var fii = []byte("fii\n001\n")

var Shapes = []Shape{
	Shape{
		[][]byte{
			[]byte{1, 1, 1, 1},
		},
	},
	Shape{
		[][]byte{
			[]byte{0, 1, 0},
			[]byte{1, 1, 1},
			[]byte{0, 1, 0},
		},
	},
	Shape{
		[][]byte{
			[]byte{0, 0, 1},
			[]byte{0, 0, 1},
			[]byte{1, 1, 1},
		},
	},
	Shape{
		[][]byte{
			[]byte{1},
			[]byte{1},
			[]byte{1},
			[]byte{1},
		},
	},
	Shape{
		[][]byte{
			[]byte{1, 1},
			[]byte{1, 1},
		},
	},
}

// Generator

type Generator[T any] struct {
	next  int
	items []T
}

func (g *Generator[T]) Next() T {
	item := g.items[g.next]
	g.next = (g.next + 1) % len(g.items)
	return item
}

// Area

type Point struct {
	x, y int
}

type Area struct {
	b            [][]byte // the first row is the top row
	falling      []Point
	highestBlock int
}

func NewArea(width, height int) Area {
	b := make([][]byte, height)
	for i := range b {
		b[i] = make([]byte, width)
		// fill b[i] with '.' runes
		for j := range b[i] {
			b[i][j] = '\x00'
		}
	}
	return Area{b, make([]Point, 0), -1}
}

/*
print the are to stdout with the falling blocks on top of the static blocks
print falling blocks as '@' and static blocks as '#'
*/
func (a *Area) Print(highestRow int) {
	for row := highestRow; row >= 0; row-- {
		for col := 0; col < len(a.b[row]); col++ {
			if a.b[row][col] == '#' {
				fmt.Print("#")
			} else {
				found := false
				for _, falling := range a.falling {
					if falling.x == col && falling.y == row {
						fmt.Print("@")
						found = true
						break
					}
				}
				if !found {
					fmt.Print(".")
				}
			}
		}
		fmt.Println()
	}
}

/*
find the new highest block in the area
*/
func (a *Area) updateHighestBlock() {
	for i, row := range a.b {
		hasNone := true
		for _, col := range row {
			if col == '#' {
				hasNone = false
				break
			}
		}
		if hasNone {
			a.highestBlock = i - 1
			break
		}
	}
}

/*
given a shape, choose the bottom left corner as the anchor point.
place the shape so that the anchor is 2 units from the left and 3 units from the highest block.
*/
func (a *Area) PlaceShape(s Shape) {
	offsetLeft := 2
	offsetHighest := 3
	highestBlock := a.highestBlock
	points := make([]Point, 0)
	rowsCnt := len(s.b)
	for ridx, row := range s.b {
		for cidx, col := range row {
			if col == 0 {
				continue
			}
			points = append(points, Point{cidx + offsetLeft, highestBlock + rowsCnt - ridx + offsetHighest})
		}
	}
	a.falling = points
}

/*
return boolean indicating if we have any falling blocks
*/
func (a *Area) HasFalling() bool {
	return len(a.falling) > 0
}

/*
freeze all falling blocks to 'b' at the correct coordinates
*/
func (a *Area) FreezeFalling() {
	for _, falling := range a.falling {
		row := falling.y
		col := falling.x
		a.b[row][col] = '#'
	}
	a.falling = make([]Point, 0)
	a.updateHighestBlock()
}

/*
first check if each falling block can move to the indicated direction (<, >, v).
if all blocks can move, move them and return true, if any block can not move, return false
and do not move any blocks.
*/
func (a *Area) MoveFalling(dir rune) bool {
	newFalling := make([]Point, 0)
	for _, falling := range a.falling {
		newFalling = append(newFalling, falling)
	}

	for i, falling := range newFalling {
		row := falling.y
		col := falling.x
		switch dir {
		case '<':
			if col == 0 {
				return false
			}
			if a.b[row][col-1] != '\x00' {
				return false
			}
			newFalling[i].x = col - 1
		case '>':
			if col == len(a.b[0])-1 {
				return false
			}
			if a.b[row][col+1] != '\x00' {
				return false
			}
			newFalling[i].x = col + 1
		case 'v':
			if row == len(a.b)-1 {
				return false
			}
			if row-1 == -1 || a.b[row-1][col] != '\x00' {
				return false
			}
			newFalling[i].y = row - 1
		}
	}
	a.falling = newFalling
	return true
}

/*
find repeating patterns from highest to lowest row
try all pattern lengths, 2..x
*/
func (a *Area) FindPattern(freshRows int) (int, int) {
	//fmt.Println("finding pattern: ", a.highestBlock, a.b[a.highestBlock], "to", checkedRow, a.b[checkedRow])
	middleIdx := a.highestBlock / 2
	//fmt.Println("finding")

	ptrs := [3]int{0, 0, 0}
	for patternLen := 10; patternLen < len(a.b)/4; patternLen++ {
		if (middleIdx-patternLen*2)-patternLen < 0 {
			continue
		}
		ptrs = [3]int{middleIdx, middleIdx - patternLen, middleIdx - (patternLen * 2)}
		nomatch := false
		for patidx := 0; patidx < patternLen; patidx++ {
			if !(SlicesEqual(a.b[ptrs[0]], a.b[ptrs[1]]) && SlicesEqual(a.b[ptrs[1]], a.b[ptrs[2]])) {
				nomatch = true
				break
			}
			//fmt.Println("found match", ptrs[0], ptrs[1], ptrs[2], "of length", patternLen)
		}
		if !nomatch {
			return ptrs[0], patternLen
		}
	}

	return 0, 0

	/*
		for row := a.highestBlock/2; row >= a.highestBlock-freshRows; row-- {
			for {
				current := string(a.b[row])
				next := string(a.b[row-1])
				if ok:=data[current]; ok != nil {
					break
				}

				for ;ok := data[current][next]; ok; {
					current = next
					next = string(a.b[row-1])
				}
			}
			for patternLen := 10; patternLen < a.highestBlock; patternLen++ {
				if row-patternLen < 0 {
					break
				}
				matches := true
				for i := 0; i < patternLen; i++ {
					for col := 0; col < len(a.b[row]); col++ {
						if a.b[row+i][col] != a.b[row-patternLen+i][col] {
							matches = false
							break
						}
					}
				}
				if matches {
					//fmt.Println("found pattern: ", patternLen, "at", row)
					return patternLen, a.highestBlock - row
				}
			}
		}
		return 0, 0*/
}

type RepeatMatcher struct {
	rowOffset   int
	patternLen  int
	shapeGenIdx int
	gasGenIdx   int
}

func play(area Area, shapeGen Generator[Shape], gasGen Generator[rune], maxblocks int64) int64 {
	repeats := map[RepeatMatcher]int{}
	startCountingAt := 0
	findRepeats := true
	simulatedHeight := 0
	prevCheckHighestBlock := 0
	var countingShapesFor *RepeatMatcher = nil

	for blockidx := int64(0); blockidx < maxblocks; blockidx++ {
		shape := shapeGen.Next()
		area.PlaceShape(shape)
		for {
			gas := gasGen.Next()
			_ = area.MoveFalling(rune(gas)) // < or >
			ok := area.MoveFalling('v')     // down
			if !ok {
				area.FreezeFalling()
				break
			}
		}
		// part 2
		if blockidx%100 == 0 {
			fmt.Println("(", maxblocks, ")", "blockidx: ", blockidx)
		}

		var matchObj *RepeatMatcher = nil
		if findRepeats {
			var matchBegins, matchLength = area.FindPattern(area.highestBlock - prevCheckHighestBlock)
			_ = matchBegins
			if matchLength > 0 {
				//fmt.Println("shape ", blockidx, "repeating block: ", matchBegins, "of length", matchLength)
				matchObj = &RepeatMatcher{(area.highestBlock - matchBegins) % matchLength, matchLength, shapeGen.next, gasGen.next}
				fmt.Println("found repeat: ", matchObj, repeats[*matchObj], "begins:", matchBegins, area.highestBlock-matchBegins, area.highestBlock)
				prevCheckHighestBlock = int(area.highestBlock)
				if _, ok := repeats[*matchObj]; !ok {
					repeats[*matchObj] = 0
				}
				repeats[*matchObj]++
			}
		}
		if matchObj != nil {
			if repeats[*matchObj] == 3 {
				countingShapesFor = matchObj
				startCountingAt = int(blockidx)
				fmt.Println("COUNTING", matchObj)
			}
			if countingShapesFor == nil {
				continue
			}
			if countingShapesFor != nil && *matchObj == *countingShapesFor && repeats[*countingShapesFor] == 4 {
				findRepeats = false
				fmt.Println("match: ", countingShapesFor, matchObj)
				canAddShapes := int(blockidx) - startCountingAt
				canAddHeight := matchObj.patternLen
				fmt.Println("found we can safely add ", canAddShapes, "shapes", canAddHeight, "height to blockIdx:", blockidx, "at height", area.highestBlock)
				for blockidx+int64(canAddShapes) < maxblocks {
					blockidx += int64(canAddShapes)
					simulatedHeight += canAddHeight
				}
				fmt.Println("now at block", blockidx)
				fmt.Println("simualted: height", simulatedHeight, "sim+highest:", simulatedHeight+area.highestBlock)
			}
		}
	}
	return int64(area.highestBlock+1) + int64(simulatedHeight)
}

/*
load 'gasPattern' as string from file using os.ReadFile from Argv[1].
then create a generator for the [] Shape.
create an Area with width 7 height 8000.
then call play() with Area and shape generator.
*/
func main() {
	gasPattern, err := os.ReadFile(os.Args[1])
	// trim gasPattern
	gasPattern = bytes.Trim(gasPattern, "\n")
	Fatal(err)
	shapeMachine := Generator[Shape]{0, Shapes}
	gasMachine := Generator[rune]{0, []rune(string(gasPattern))}
	area := NewArea(7, 2022*4)
	p1Result := play(area, shapeMachine, gasMachine, int64(2022))
	fmt.Printf("p1: %d\n", p1Result)

	/*
		shapeMachine = Generator[Shape]{0, Shapes}
		gasMachine = Generator[rune]{0, []rune(string(gasPattern))}
		area = NewArea(7, 1000000)
		p2Result := play(area, shapeMachine, gasMachine, int64(1000000000000))
		fmt.Printf("p2: %d\n", p2Result)
	*/
}
