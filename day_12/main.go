package main

import (
	"fmt"
	"os"
	"strings"
)

func Fatal(err error) {
	if err != nil {
		panic(err)
	}
}

type Tuple[T int] struct {
	x, y T
}

type Tile struct {
	isStart      bool
	isEnd        bool
	letter       rune
	shortestPath int
}

/*
parse a map of the form:
Sabqponm
abcryxxl
accszExk
acctuvwj
abdefghi

S is the start point. E is the end point. each letter a-z is a height, a is the lowest, z is the highest.
S is always at a and E is always at z.
I can move up or down one letter at a time, but not diagonally.
I can move to a tile maximum 1 larger than my current one. I can always move to all lower tiles.
*/
func ParseMap(input string) [][]Tile {
	var tiles [][]Tile
	for _, line := range strings.Split(input, "\n") {
		// skip empty lines
		if len(line) == 0 {
			continue
		}
		var row []Tile
		for _, letter := range line {
			tile := Tile{letter: letter, shortestPath: -1, isStart: letter == 'S', isEnd: letter == 'E'}
			if tile.isStart {
				tile.letter = 'a'
			}
			if tile.isEnd {
				tile.letter = 'z'
			}
			row = append(row, tile)
		}
		tiles = append(tiles, row)
	}
	return tiles
}

var counter int = 0

/*
find the shortest path from 'at' to E.
do this recursively.
tiles is [y][x]
*/
func FindShortestPath(from Tuple[int], at Tuple[int], stepsTaken int, tiles [][]Tile) int {
	// if we've walked off the edge of the map, we're done
	if at.x < 0 || at.y < 0 || at.x >= len(tiles[0]) || at.y >= len(tiles) {
		return -1
	}

	// if we've walked onto a tile that's too high, we're done
	if tiles[at.y][at.x].letter > tiles[from.y][from.x].letter+1 {
		return -1
	}

	// if we've already been here and found a shorter path, we're done
	if tiles[at.y][at.x].shortestPath > 0 && tiles[at.y][at.x].shortestPath <= stepsTaken {
		return -1
	}

	// if we've found a path to E that's shorter than the one we're on, use that one
	if tiles[at.y][at.x].isEnd && (tiles[at.y][at.x].shortestPath == -1 || tiles[at.y][at.x].shortestPath > stepsTaken) {
		tiles[at.y][at.x].shortestPath = stepsTaken
		return stepsTaken
	}

	// we've walked onto a tile that's the right height. walk onto it.
	tiles[at.y][at.x].shortestPath = stepsTaken

	coordinates := []Tuple[int]{
		{at.x - 1, at.y},
		{at.x + 1, at.y},
		{at.x, at.y - 1},
		{at.x, at.y + 1},
	}

	shortest := -1
	for _, coordinate := range coordinates {
		shortestPath := FindShortestPath(at, coordinate, stepsTaken+1, tiles)
		if shortestPath > 0 && (shortest == -1 || shortestPath < shortest) {
			shortest = shortestPath
		}
	}

	return shortest
}

/*
read input from file at os.ReadFile(os.Args[1])
parse the file and then find the shortest path from S to E
*/
func main() {
	input, err := os.ReadFile(os.Args[1])
	Fatal(err)
	tiles := ParseMap(string(input))
	// find start tile
	var start Tuple[int]
	for y, row := range tiles {
		for x, tile := range row {
			if tile.isStart {
				start = Tuple[int]{x, y}
			}
		}
	}
	// print starting location
	shortest := FindShortestPath(start, start, 0, tiles)
	fmt.Println("Part1:", shortest)

	// find each 'a' and find the shortest path from there to E
	shortestAPath := -1
	for y, row := range tiles {
		for x, tile := range row {
			if tile.letter == 'a' {
				shortest := FindShortestPath(Tuple[int]{x, y}, Tuple[int]{x, y}, 0, tiles)
				if shortest > 0 && (shortestAPath == -1 || shortest < shortestAPath) {
					shortestAPath = shortest
				}
			}
		}
	}
	fmt.Println("Part2:", shortestAPath)
}
