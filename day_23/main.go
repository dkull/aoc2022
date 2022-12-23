package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

/*
Utils*
*/
func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

/*
Structures
*/

type Vector struct {
	Name string
	X, Y int
}

// the ordering is important, at least for the first 4
var Vectors = []Vector{
	{"North", 0, -1},
	{"South", 0, 1},
	{"West", -1, 0},
	{"East", 1, 0},
	{"NorthEast", 1, -1},
	{"NorthWest", -1, -1},
	{"SouthEast", 1, 1},
	{"SouthWest", -1, 1},
}

type Point struct {
	x int
	y int
}

func (p Point) AddVector(v Vector) Point {
	return Point{p.x + v.X, p.y + v.Y}
}

func (p *Point) Get8Neighbors() []Point {
	neighbors := []Point{}
	for _, v := range Vectors {
		neighbors = append(neighbors, p.AddVector(v))
	}
	return neighbors
}

/*
gfiven a point and a vector, give me not just the point, but the
two neighbors too, basically replace the 0 vec with -1 and 1 too.
*/
func (p Point) GetLookPoints(v Vector) []Point {
	var out []Point
	if v.X == 0 {
		out = append(out, Point{p.x - 1, p.y + v.Y})
		out = append(out, Point{p.x, p.y + v.Y})
		out = append(out, Point{p.x + 1, p.y + v.Y})
	} else {
		out = append(out, Point{p.x + v.X, p.y - 1})
		out = append(out, Point{p.x + v.X, p.y})
		out = append(out, Point{p.x + v.X, p.y + 1})
	}
	return out
}

type Elf struct {
	Name          string
	Position      Point
	Directions    []Vector
	DirectionsIdx int
}

func (e *Elf) HaveAnyNeighbor(otherElves map[Point]*Elf) bool {
	neighbors := e.Position.Get8Neighbors()
	for _, neighbor := range neighbors {
		if _, ok := otherElves[neighbor]; ok {
			return true
		}
	}
	return false
}

func (e *Elf) ProposeMove(otherElves map[Point]*Elf) *Point {
	var result *Point
dirLoop:
	for x := 0; x < len(e.Directions); x++ {
		dirIdx := (e.DirectionsIdx + x) % len(e.Directions)
		lookPoints := e.Position.GetLookPoints(e.Directions[dirIdx])
		// if any elf occupying the 3 points in the direction of the vec
		// then we can't move there
		for _, lookPoint := range lookPoints {
			if _, ok := otherElves[lookPoint]; ok {
				fmt.Println("elf", e.Name, "can't move", e.Directions[dirIdx], "because of elf at", lookPoint)
				continue dirLoop
			}
		}
		fmt.Println("elf", e.Name, "can move", e.Directions[dirIdx])
		// the first one consideredd is move to lat
		result = &lookPoints[1]
		break
	}
	e.DirectionsIdx = (e.DirectionsIdx + 1) % len(e.Directions)
	return result
}

/*
Functions
*/

/*
Parse an input of a map of the form:
....#..
..###.#
#...#.#
.#...##
#.###..
##.#.##
.#..#..

. = ground
# = elf
*/
func ParseMap(input string) (map[Point]bool, []Elf) {
	elves := []Elf{}
	ground := map[Point]bool{}
	lines := strings.Split(input, "\n")
	for y, line := range lines {
		for x, char := range line {
			if char == '#' {
				elves = append(elves, Elf{
					Name:          fmt.Sprintf("%d", len(elves)),
					Position:      Point{x: x, y: y},
					Directions:    Vectors[0:4], // use only N,S,E,W
					DirectionsIdx: 0,
				})
			} else if char == '.' {
				ground[Point{x, y}] = true
			}
		}
	}
	return ground, elves
}

/*
Find the bounding box around all elves, then draw the map
with '.' for ground and '#' for each elf.
*/
func DrawMap(elves map[Point]*Elf) {
	minX, minY := 1000, 1000
	maxX, maxY := 0, 0
	for point := range elves {
		if point.x < minX {
			minX = point.x
		}
		if point.y < minY {
			minY = point.y
		}
		if point.x > maxX {
			maxX = point.x
		}
		if point.y > maxY {
			maxY = point.y
		}
	}
	fmt.Println("TL corner:", minX, minY)
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if elf, ok := elves[Point{x, y}]; ok {
				if len(elf.Name) == 1 {
					fmt.Print(elf.Name)
				} else {
					fmt.Print("#")
				}
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

/*
find bounding box and count all '.' cells
*/
func CalcScore(elves []Elf) int {
	minX, minY := math.MaxInt, math.MaxInt
	maxX, maxY := math.MinInt, math.MinInt
	for _, elf := range elves {
		if elf.Position.x < minX {
			minX = elf.Position.x
		}
		if elf.Position.y < minY {
			minY = elf.Position.y
		}
		if elf.Position.x > maxX {
			maxX = elf.Position.x
		}
		if elf.Position.y > maxY {
			maxY = elf.Position.y
		}
	}
	elfMap := map[Point]*Elf{}
	for _, elf := range elves {
		elfMap[elf.Position] = &elf
	}
	// count '.' in the bounding box
	score := 0
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if _, ok := elfMap[Point{x, y}]; !ok {
				score++
			}
		}
	}
	return score
}

func Part1(elves []Elf) int {
	for round := 1; round <= 10; round++ {
		fmt.Println("\n==== Round", round, "====\n")
		elfAt := map[Point]*Elf{}
		for i := range elves {
			elf := &elves[i]
			elfAt[elf.Position] = elf
		}

		proposedMoves := map[Point][]*Elf{}
		for _, elf := range elfAt {
			// for each neighbor, if
			haveNeighbor := elf.HaveAnyNeighbor(elfAt)
			// we do the propose even if no neighbors, becaues
			// we always need to rotate the direction
			proposedMove := elf.ProposeMove(elfAt)
			if !haveNeighbor {
				fmt.Println("Elf", elf.Name, "has no neighbor")
				continue
			}
			if proposedMove == nil {
				fmt.Println("Elf", elf.Name, "has no proposed move")
				continue
			}

			// queue up the proposed move
			if _, ok := proposedMoves[*proposedMove]; !ok {
				proposedMoves[*proposedMove] = []*Elf{}
			}
			proposedMoves[*proposedMove] = append(proposedMoves[*proposedMove], elf)
			//fmt.Println("elf", elf.Name, "proposes to move to", proposedMove)
		}

		// only do moves to points where only 1 elf is going
		for point, elfList := range proposedMoves {
			if len(elfList) == 1 {
				fmt.Println("Elf", elfList[0].Position, "moves to", point)
				elfList[0].Position = point
			}
		}
		// Map drawing
		elfAt = map[Point]*Elf{}
		for i := range elves {
			elf := &elves[i]
			elfAt[elf.Position] = elf
		}
		DrawMap(elfAt)
		// END map drawing

	}
	return CalcScore(elves)
}

/*
Main
*/
func main() {
	data, err := os.ReadFile(os.Args[1])
	Fatal(err)
	_, elves := ParseMap(string(data))
	for elfIdx, elf := range elves {
		fmt.Printf("Elf %d: %v\n", elfIdx, elf)
	}
	fmt.Println("Part 1:", Part1(elves))
}
