package main

import (
	"bufio"
	"fmt"
	"os"
)

func Fatal(err error) {
	if err != nil {
		panic(err)
	}
}

func assertEqual(expected, actual interface{}) {
	if expected != actual {
		panic("expected " + fmt.Sprint(expected) + " but got " + fmt.Sprint(actual))
	}
}

/*
a point defined by x,y
the point has a 'treeHeight' number
the point has a 'visibleCount' number
the point has a 'scenicScore' number
*/
type Point struct {
	x, y, treeHeight, visibleCount, scenicScore int
}

/*
read in a file and add it into a 2d array
where each character is mapped to a Point
x and y are the coordinates of the point
visibleCount is 0 by default
treeHeight is the character ('0'-'9') converted to int using Atoi
scenicScore is 1 by default
*/
func readInFile(filename string) [][]Point {
	file, err := os.Open(filename)
	Fatal(err)
	defer file.Close()

	var points [][]Point

	scanner := bufio.NewScanner(file)
	for y := 0; scanner.Scan(); y++ {
		var row []Point
		for x, char := range scanner.Text() {
			row = append(row, Point{x, y, int(char - '0'), 0, 1})
		}
		points = append(points, row)
	}
	return points
}

/*
take in a 2d array of points and return the
dimensions
*/
func getDimensions(points [][]Point) (int, int) {
	return len(points[0]), len(points)
}

/*
transpose a 2d array of points 90 degrees to the right
assume the array is square
*/
func transposeRight(points [][]Point) [][]Point {
	var result [][]Point
	for i := 0; i < len(points); i++ {
		var row []Point
		for j := len(points) - 1; j >= 0; j-- {
			row = append(row, points[j][i])
		}
		result = append(result, row)
	}
	return result
}

/*
given a 2d array of points, go over each point and check each element to the
right of it. if our treeHeight <= other treeHeight the we are blocked from that
direction and we should continue to next x and not increment visibleCount.
take a pointer to our current point.
increment its visibleCount before checking the rest of the points to the right.
if any tree is not shorter than it, decrement the visibleCount and move on to the next point.
keep a count of each tree to our right that we considered.
multiply the points scenicScore by the number of trees we considered.
return the modified points
*/
func checkVisibility(points [][]Point) [][]Point {
	for y, row := range points {
		for x, point := range row {
			point.visibleCount++
			count := 0
			for i := x + 1; i < len(row); i++ {
				count++
				if point.treeHeight <= row[i].treeHeight {
					point.visibleCount--
					break
				}
			}
			point.scenicScore *= count
			points[y][x] = point
		}
	}
	return points
}

/*
take in a 2d array of points.
transpose the array four times and each times calling checkVisibility on it
return the points
*/
func checkAllDirections(points [][]Point) [][]Point {
	for i := 0; i < 4; i++ {
		checkVisibility(points)
		points = transposeRight(points)
	}
	return points
}

/*
count all points that have a visibleCount of at least 1
*/
func countVisible(points [][]Point) int {
	count := 0
	for _, row := range points {
		for _, point := range row {
			if point.visibleCount >= 1 {
				count++
			}
		}
	}
	return count
}

/*
find the highest scenicSCore in the array
*/
func findHighestScenicScore(points [][]Point) int {
	highest := 0
	for _, row := range points {
		for _, point := range row {
			if point.scenicScore > highest {
				highest = point.scenicScore
			}
		}
	}
	return highest
}

/*
assign Argc[1] to filename
call readInFile and getDimensions
assert the dimensions are equal
call checkAllDirections and countVisible
print the result as Part1: <result>
print the result of findHighestScenicScore as Part2: <result>
*/
func main() {
	filename := os.Args[1]
	points := readInFile(filename)
	width, height := getDimensions(points)
	assertEqual(width, height)
	points = checkAllDirections(points)
	fmt.Println("Part1:", countVisible(points))
	fmt.Println("Part2:", findHighestScenicScore(points))
}
