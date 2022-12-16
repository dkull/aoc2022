package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// utility

func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func ManhattanDistance(p1, p2 Point) int {
	return Abs(p1.X-p2.X) + Abs(p1.Y-p2.Y)
}

// task

type Point struct {
	X, Y int
}

type Fact struct {
	Sensor Point
	Beacon Point
}

/*
parse a line of the format, extract the two number pairs to Point:
> Sensor at x=20, y=14: closest beacon is at x=25, y=17
Then collect the two points into a Fact.
*/
func ParseLine(line string) Fact {
	var p1, p2 Point
	_, err := fmt.Sscanf(line, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &p1.X, &p1.Y, &p2.X, &p2.Y)
	Fatal(err)
	return Fact{p1, p2}
}

/*
read a file using os.ReadFile from Args[1].
parse each line using ParseLine into a Fact.
return a list of facts.
*/
func ReadFacts() []Fact {
	file, err := os.ReadFile(os.Args[1])
	Fatal(err)
	lines := strings.Split(string(file), "\n")
	facts := make([]Fact, len(lines)-1) // last line is empty
	for i, line := range lines {
		if line == "" {
			continue
		}
		facts[i] = ParseLine(line)
	}
	return facts
}

/*
given a list of facts, return the leftmost point by manhattan distance
from the leftmost Sensor. then move that point even more left by the
manhattan distance between the Sensor and Beacon.
*/
func LeftmostPOI(facts []Fact, furthest int) Point {
	var leftmostFact Fact = facts[0]
	for _, fact := range facts {
		if fact.Sensor.X < leftmostFact.Sensor.X {
			leftmostFact = fact
		}
	}
	return Point{leftmostFact.Sensor.X - furthest, leftmostFact.Sensor.Y}
}

func RightmostPOI(facts []Fact, furthest int) Point {
	var rightmostFact Fact = facts[0]
	for _, fact := range facts {
		if fact.Sensor.X > rightmostFact.Sensor.X {
			rightmostFact = fact
		}
	}
	return Point{rightmostFact.Sensor.X + furthest, rightmostFact.Sensor.Y}
}

/*
find the manhattan distance of the sensor and beacon that are the farthest apart.
*/
func FurthestDistance(facts []Fact) int {
	var maxDistance int
	for _, fact := range facts {
		distance := ManhattanDistance(fact.Sensor, fact.Beacon)
		if distance > maxDistance {
			maxDistance = distance
		}
	}
	return maxDistance
}

/*
given a list of facts, a Y coordinate and a left and right X coordinate,
find for each position from (leftX, Y) to (rightX, Y) if this point is
closer to a sensor than that sensors is to its beacon. if so, count that.
*/
func CountClean(facts []Fact, Y, leftX, rightX int) int {
	count := 0
	for x := leftX; x <= rightX; x++ {
		for _, fact := range facts {
			// don't count the beacons we hit
			if fact.Beacon.X == x && fact.Beacon.Y == Y {
				continue
			}
			if fact.Sensor.X == x && fact.Sensor.Y == Y {
				continue
			}
			if ManhattanDistance(fact.Sensor, Point{x, Y}) <= ManhattanDistance(fact.Sensor, fact.Beacon) {
				// print current progress of x and count
				count++
				break
			}
		}
	}
	return count
}

/*
find if a given point is not closer to any Sensor than
that sensors beacon is to the sensor
*/
func IsInvisible(facts []Fact, point Point) bool {
	for _, fact := range facts {
		sensorToPoint := ManhattanDistance(fact.Sensor, point)
		sensorToBeacon := ManhattanDistance(fact.Sensor, fact.Beacon)
		if sensorToPoint <= sensorToBeacon {
			return false
		}
	}
	return true
}

/*
determine if two facts overlap
*/
func Overlap(factA, factB Fact) bool {
	adist := ManhattanDistance(factA.Sensor, factA.Beacon)
	bdist := ManhattanDistance(factB.Sensor, factB.Beacon)
	return ManhattanDistance(factA.Sensor, factB.Sensor) <= adist+bdist
}

func AlmostThouching(factA, factB Fact) bool {
	adist := ManhattanDistance(factA.Sensor, factA.Beacon)
	bdist := ManhattanDistance(factB.Sensor, factB.Beacon)
	return ManhattanDistance(factA.Sensor, factB.Sensor) == adist+bdist+2
}

/*
given a list of facts, do a matrix comparison between all of them.
find factA and  factB that overlap. and factC and factD that overlap.
and also factC overlaps both factA and factB. and factD overlaps both
factA and factB. this sets up the rects in a suitable way for our final step.
*/
func FindFactPair(facts []Fact) (Fact, Fact, Fact, Fact) {
	for _, A := range facts {
		for _, B := range facts {
			for _, C := range facts {
				for _, D := range facts {
					if A.Sensor.X > B.Sensor.X || A.Sensor.Y > B.Sensor.Y {
						continue
					}
					if C.Sensor.X < D.Sensor.X || C.Sensor.Y > D.Sensor.Y {
						continue
					}
					if !(AlmostThouching(A, B) && AlmostThouching(C, D)) {
						continue
					}
					if !(Overlap(C, A) && Overlap(C, B)) {
						continue
					}
					if !(Overlap(D, A) && Overlap(D, B)) {
						continue
					}
					return A, B, C, D
				}
			}
		}
	}
	panic(fmt.Sprintf("no fact pair found: %v", facts))
}

// 4736899 is low
// 4347487 is low
// 4347486 is low
func main() {
	facts := ReadFacts()
	furthest := FurthestDistance(facts)
	leftmost := LeftmostPOI(facts, furthest)
	rightmost := RightmostPOI(facts, furthest)
	countP1 := CountClean(facts, 2000000, leftmost.X, rightmost.X)
	fmt.Printf("Part1: %d\n", countP1)

	// part 2
	_, _, _, factD := FindFactPair(facts)

	// sensor D manhattan distance to its sensor
	manD := ManhattanDistance(factD.Sensor, factD.Beacon)
	// move to highest point on D, then +1
	topD := Point{factD.Sensor.X, factD.Sensor.Y - manD - 1}
	// start moving down-right from topD and check IsInvisible
	// until we find a point that is invisible
	for {
		if IsInvisible(facts, topD) {
			if topD.X >= 0 && topD.X <= 4000000 && topD.Y >= 0 && topD.Y <= 4000000 {
				break
			} else {
				break
			}
		}
		topD.X++
		topD.Y++
	}
	// multiply X by 4000000 and add Y
	fmt.Printf("Part2: %d\n", topD.X*4000000+topD.Y)
}
