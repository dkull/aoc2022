package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

/*
holds a direction and a distance
*/
type Direction int

const (
	UP Direction = iota
	DOWN
	LEFT
	RIGHT
)

func Fatal(err error) {
	if err != nil {
		panic(err)
	}
}

// holds a direction and a distance
type Move struct {
	direction Direction
	distance  int
}

/*
read a file from argument 'filename' and split it into lines
remove the last empty line
*/
func readLines(filename string) []string {
	file, err := os.Open(filename)
	Fatal(err)
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	Fatal(scanner.Err())
	return lines
}

/*
moves are lines of the form "R 4" or "D 21"
the direction is the first thing, and the distance is the second thing
there is a space in the middle
parse the number using strconv.Atoi
*/
func parseMoves(lines []string) []Move {
	var moves []Move
	for _, line := range lines {
		var move Move
		switch line[0] {
		case 'U':
			move.direction = UP
		case 'D':
			move.direction = DOWN
		case 'L':
			move.direction = LEFT
		case 'R':
			move.direction = RIGHT
		}
		distance, err := strconv.Atoi(line[2:])
		Fatal(err)
		move.distance = distance
		moves = append(moves, move)
	}
	return moves
}

/*
calculate all coordinates that form a square of side length sideLen.
center the square around the point X, Y
XXXXX
X...X
X.S.X
X...X
XXXXX
Copilot didn't manage to solve this, ChatGPT did (with some extra input from me).
*/
func neighbors(x, y, sideLen int) [][2]int {
	var neighbors [][2]int
	for i := x - sideLen/2; i <= x+sideLen/2; i++ {
		for j := y - sideLen/2; j <= y+sideLen/2; j++ {
			if i == x-sideLen/2 || i == x+sideLen/2 || j == y-sideLen/2 || j == y+sideLen/2 {
				neighbors = append(neighbors, [2]int{i, j})
			}
		}
	}
	return neighbors
}

/*
the arguments are a list of 'moves' and the number of moving things in total.
the 'track' argument signifies which(nth) objects unique visited locations count
we return as the result. collect the unique locations visited by that
item locations in a map.

allocate all the moving things as coordinates into an array.
loop over the moves, but only the first item will move like the 'move' says.
then move all following moving things, by following these rules. the first one does not do this:
 1. call neighbors(my_x, my_y, 5)
 2. check if the previous item coordinate is one of those
 3. if they do set my current position as the result of findNextPosition()
*/
func runMoves(moves []Move, tailsCount int, track int) int {
	// allocate all the moving things as coordinates into an array.
	var tails [][2]int
	for i := 0; i < tailsCount; i++ {
		tails = append(tails, [2]int{0, 0})
	}

	// collect the unique locations visited by that item locations in a map.
	visited := make(map[[2]int]bool)
	visited[tails[track]] = true

	// loop over the moves, but only the first item will move like the 'move' says.
	// then move all following moving things, by following these rules. the first one does not do this:
	for _, move := range moves {
		// 1. call neighbors(my_x, my_y, 5)
		// 2. check if the previous item coordinate is one of those
		// 3. if they do set my current position as the result of findNextPosition()
		for step := 0; step < move.distance; step++ {
			for i := 0; i < tailsCount; i++ {
				if i == 0 {
					switch move.direction {
					case UP:
						tails[i][1] += 1
					case DOWN:
						tails[i][1] -= 1
					case LEFT:
						tails[i][0] -= 1
					case RIGHT:
						tails[i][0] += 1
					}
				} else {
					neighbors := neighbors(tails[i][0], tails[i][1], 5)
					for _, neighbor := range neighbors {
						if neighbor == tails[i-1] {
							tails[i] = findNextPosition(tails[i-1], tails[i])
							break
						}
					}
				}
				if i == track {
					visited[tails[i]] = true
				}
			}
		}
	}
	return len(visited)
}

/*
if the head and tail are in the same column or row, move the tail closer to the head by 1.
if they are not in the same row or column, move the tail closer to the head by 1, but diagonally.
*/
func findNextPosition(head, tail [2]int) [2]int {
	if head[0] == tail[0] {
		if head[1] > tail[1] {
			tail[1]++
		} else {
			tail[1]--
		}
	} else if head[1] == tail[1] {
		if head[0] > tail[0] {
			tail[0]++
		} else {
			tail[0]--
		}
	} else {
		if head[0] > tail[0] {
			tail[0]++
		} else {
			tail[0]--
		}
		if head[1] > tail[1] {
			tail[1]++
		} else {
			tail[1]--
		}
	}
	return tail
}

/*
read file from os.Args[1] using readLines. parse the moves using parseMoves().
run the moves with 2,1 and print the result as Part1:
run the moves with 10,9 and print the result as Part2:
*/
func main() {
	lines := readLines(os.Args[1])
	moves := parseMoves(lines)
	fmt.Println("Part1:", runMoves(moves, 2, 1))
	fmt.Println("Part2:", runMoves(moves, 10, 9))
}
