package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Point struct {
	X, Y int
}

type Blizzard struct {
	Pos Point
	Dir byte
}

type Mover struct {
	Pos   Point
	Moves int
}

type Map struct {
	Columns   int
	Lines     int
	Teleports map[Point]Point
	Finish    Point
}

/*
Parse a map of the form:
#.######
#>>.<^<#
#.<..<<#
#>v.><>#
#<^v^^>#
######.#

'#' are walls, '>', '<', '^' and 'v' are blizzards, which move
where they point. Entry point is the empty square in the first row,
exit is the empty square in the last row.
*/
func Parse(input string) (Point, []Blizzard, Map) {
	teleport := make(map[Point]Point)
	blizzards := make([]Blizzard, 0)
	start := Point{}
	finish := Point{}

	lines := strings.Split(input, "\n")
	// remove possible last empty line
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	for y, line := range lines {
		// trim whitespace from line
		line = strings.TrimSpace(line)
		for x, c := range line {
			// teleports
			// buggy, most edge boxes become teleporters
			if y == 0 && c == '#' {
				teleport[Point{x, 0}] = Point{x, len(lines) - 2}
				teleport[Point{x, len(lines) - 1}] = Point{x, 1}
			}
			if x == 0 && c == '#' {
				teleport[Point{0, y}] = Point{len(line) - 2, y}
				teleport[Point{len(line) - 1, y}] = Point{1, y}
			}

			if y == 0 && c == '.' {
				start = Point{x, 0}
			} else if y == len(lines)-1 && c == '.' {
				finish = Point{x, len(lines) - 1}
			} else {
				if c == '>' || c == '<' || c == '^' || c == 'v' {
					blizzards = append(blizzards, Blizzard{Point{x, y}, byte(c)})
				}
			}
		}
	}

	mapp := Map{
		Lines:     len(lines),
		Columns:   len(lines[0]) - 1, // newline
		Teleports: teleport,
		Finish:    finish,
	}
	if mapp.Finish.X == 0 && mapp.Finish.Y == 0 {
		panic("No finish found")
	}
	return start, blizzards, mapp
}

/*
move each blizzard towards its direction, if it ends up in
a teleport tile, teleport it.
*/
func moveBlizzards(bliz []Blizzard, teleport map[Point]Point) []Blizzard {
	for i, b := range bliz {
		switch b.Dir {
		case '>':
			bliz[i] = Blizzard{Point{b.Pos.X + 1, b.Pos.Y}, b.Dir}
		case '<':
			bliz[i] = Blizzard{Point{b.Pos.X - 1, b.Pos.Y}, b.Dir}
		case '^':
			bliz[i] = Blizzard{Point{b.Pos.X, b.Pos.Y - 1}, b.Dir}
		case 'v':
			bliz[i] = Blizzard{Point{b.Pos.X, b.Pos.Y + 1}, b.Dir}
		}
		// teleport blizzard
		if t, ok := teleport[bliz[i].Pos]; ok {
			bliz[i] = Blizzard{t, bliz[i].Dir}
		}
	}
	return bliz
}

func run2(start Point, bliz []Blizzard, mapp Map, targets []Point) int {
	movers := make(map[Point]Mover)
	movers[start] = Mover{start, 0}

	mapp.Finish = targets[0]
	targets = targets[1:]

	// for loop while movers[mapp.Final] does not exist
	// we move the blizzards, then create new movers in
	// every direction around all movers, then we
	// delete the movers that can't exist
	// end when one mover reaches the Finish
	for {
		// move all blizzards
		bliz = moveBlizzards(bliz, mapp.Teleports)
		// create a blizzard lookup
		blizMap := make(map[Point]bool)
		for _, b := range bliz {
			blizMap[b.Pos] = true
		}
		// create new movers
		newMovers := make(map[Point]Mover)
	moversFor:
		for _, m := range movers {
			// mover is clear, propagate
			newMoverPoints := []Point{
				{m.Pos.X + 1, m.Pos.Y},
				{m.Pos.X - 1, m.Pos.Y},
				{m.Pos.X, m.Pos.Y + 1},
				{m.Pos.X, m.Pos.Y - 1},
				{m.Pos.X, m.Pos.Y},
			}
			for _, newMoverPoint := range newMoverPoints {
				// check if new mover is on Finish
				if newMoverPoint == mapp.Finish {
					if len(targets) == 0 {
						return m.Moves + 1
					} else {
						newMovers = make(map[Point]Mover)
						newMovers[newMoverPoint] = Mover{newMoverPoint, m.Moves + 1}
						// we have a bug in the teleporter creator
						// where target boxes are teleporters
						delete(mapp.Teleports, newMoverPoint)
						delete(mapp.Teleports, targets[0])
						mapp.Finish = targets[0]
						targets = targets[1:]
						break moversFor
					}
				}
				if _, ok := newMovers[newMoverPoint]; ok {
					continue
				}
				// check if new mover is in a blizzard
				if blizMap[newMoverPoint] {
					continue
				}
				// check if new mover is in a teleporter, continue
				if _, ok := mapp.Teleports[newMoverPoint]; ok {
					continue
				}
				// check if new mover is out of bounds top or bottom
				if newMoverPoint.Y == -1 || newMoverPoint.Y == mapp.Lines {
					continue
				}
				newMovers[newMoverPoint] = Mover{newMoverPoint, m.Moves + 1}
			}
		}
		movers = newMovers
	}
}

/*
504 too high
458 too high
442 is too high
*/
func main() {
	// read file from Args[1] using os.ReadFile
	data, err := os.ReadFile(os.Args[1])
	Fatal(err)
	start, blizzards, mapp := Parse(string(data))

	// part2:
	targets := []Point{
		{mapp.Columns - 2, mapp.Lines - 1},
		{1, 0},
		{mapp.Columns - 2, mapp.Lines - 1},
	}
	// for part1, remove the 2 last targets

	moves := run2(start, blizzards, mapp, targets)
	if moves < math.MaxInt {
		fmt.Println("Answer:", moves)
	}
}
