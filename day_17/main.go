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
			b[i][j] = '.'
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
	fmt.Println("highest block:", a.highestBlock)
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
	fmt.Println("freezing falling blocks")
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
			if a.b[row][col-1] != '.' {
				return false
			}
			newFalling[i].x = col - 1
		case '>':
			if col == len(a.b[0])-1 {
				return false
			}
			if a.b[row][col+1] != '.' {
				return false
			}
			newFalling[i].x = col + 1
		case 'v':
			if row == len(a.b)-1 {
				return false
			}
			if row-1 == -1 || a.b[row-1][col] != '.' {
				return false
			}
			newFalling[i].y = row - 1
		}
	}
	a.falling = newFalling
	return true
}

func play(area Area, shapeGen Generator[Shape], gasGen Generator[rune], maxblocks int) int {
	for blockidx := 0; blockidx < maxblocks; blockidx++ {
		fmt.Printf("=== round %d ===\n", blockidx)
		shape := shapeGen.Next()
		area.PlaceShape(shape)
		//area.Print(20)
		for {
			gas := gasGen.Next()
			_ = area.MoveFalling(rune(gas)) // < or >
			ok := area.MoveFalling('v')     // down
			if !ok {
				area.FreezeFalling()
			}
			if !area.HasFalling() {
				break
			}
		}
	}
	return area.highestBlock
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
	area := NewArea(7, 9000)
	p1Result := play(area, shapeMachine, gasMachine, 2022)
	fmt.Printf("p1: %d\n", p1Result+1) // +1 because we start at 0
}
