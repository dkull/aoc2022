package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// utility

func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Abs[T int](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func ToUnit(x int) int {
	if x == 0 {
		return 0
	}
	return x / Abs(x)
}

// task

type Pos struct {
	X, Y int
}

func (p *Pos) Add(pos Pos) {
	p.X += pos.X
	p.Y += pos.Y
}

/*
return all Pos points between two Pos points.
the two input points are guaranteed to be on the same row or column
first figure out the delta of x and y, then move from point to point
Handwritten
*/
func GetConnectingPoints(a, b Pos) []Pos {
	deltaX := b.X - a.X
	deltaY := b.Y - a.Y
	if deltaX != 0 && deltaY != 0 {
		panic("points are not on the same row or column")
	}
	// calculate the unit vector Pos, deltas can be 0
	unitVector := Pos{ToUnit(deltaX), ToUnit(deltaY)}

	// start moving by the unit vector from point a to point b
	// saving all the points on the way
	var points []Pos = []Pos{a, b}
	a.Add(unitVector)
	for p := a; p != b; p.Add(unitVector) {
		points = append(points, p)
	}
	return points
}

type PlayField struct {
	emitter      Pos
	atRest       map[Pos]int
	atRestBuried map[Pos]int
	stone        map[Pos]int
	falling      map[Pos]int
	lowestStoneY int
	floor        bool
	hitAbyss     bool
}

/*
the top is determined by the highest emitter, left right and bottom by the stones
draw emitters as '*', stone as '#' falling as '.' and atRest as 'o'.
small y is up, larger y is down.
*/
func (p *PlayField) Draw() {
	padding := 10
	// find the top, left, right and bottom
	// initialize them all to the first emitter coordinates
	emitter := p.emitter
	top := emitter.Y - padding
	left := emitter.X - padding
	right := emitter.X + padding
	bottom := emitter.Y + padding
	if emitter.Y < top {
		top = emitter.Y
	}
	if emitter.X < left {
		left = emitter.X
	}
	if emitter.X > right {
		right = emitter.X
	}
	for pos, _ := range p.stone {
		if pos.Y < top {
			top = pos.Y
		}
		if pos.Y > bottom {
			bottom = pos.Y
		}
		if pos.X < left {
			left = pos.X
		}
		if pos.X > right {
			right = pos.X
		}
	}
	// draw
	for y := top; y <= bottom; y++ {
		for x := left; x <= right; x++ {
			pos := Pos{x, y}
			if _, ok := p.atRest[pos]; ok {
				fmt.Print("o")
			} else if _, ok := p.atRestBuried[pos]; ok {
				fmt.Print("X")
			} else if _, ok := p.falling[pos]; ok {
				fmt.Print(".")
			} else if _, ok := p.stone[pos]; ok {
				fmt.Print("#")
			} else if p.emitter == pos {
				fmt.Print("*")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

/*
stone is a list of points that are stones, and the straight
lines between the points are stones.
for example:

	x1,y1->x2,y2->x3,y3->x4,y4

split the line by "->" and then split each point by ",".
parse the coordinates into a list of Pos points.
then iterate over pairs of points and find the connecting points using GetConnectingPoints.
add all the connecting points to the stone map.
*/
func (p *PlayField) ParseStoneFromLine(line string) {
	// split the line by "->"
	parts := strings.Split(line, " -> ")
	// parse the coordinates into a list of Pos points
	var points []Pos
	for _, part := range parts {
		// split each point by ","
		coords := strings.Split(part, ",")
		x, err := strconv.Atoi(coords[0])
		Fatal(err)
		y, err := strconv.Atoi(coords[1])
		Fatal(err)
		points = append(points, Pos{x, y})
	}
	// iterate over pairs of points and find the connecting points
	for i := 0; i < len(points)-1; i++ {
		for _, point := range GetConnectingPoints(points[i], points[i+1]) {
			// we count the end points multiple times, but that's ok
			// since we are using a map
			p.stone[point] = 0
			if point.Y > p.lowestStoneY {
				p.lowestStoneY = point.Y
			}
		}
	}
}

/*
add new falling objects to the falling map.
starting from each emitter.
*/
func (p *PlayField) AddNewFalling() {
	p.falling[p.emitter] = 0
}

/*
move all falling objects one step down.
atRest and stone objects block the way.
if the object can move down move down.
if not, try to move down-left, then down-right.
remove from falling and add to atRest if cannot move any more.
if a falling object falls lower than the lowest stone, mark the hitAbyss.
*/
func (p *PlayField) MoveFalling() {
	for pos, _ := range p.falling {
		// check which coordinates are free in atRest and stone
		var moved bool = false
		for _, delta := range []Pos{{0, 1}, {-1, 1}, {1, 1}} {
			newPos := pos
			newPos.Add(delta)
			if _, ok := p.atRest[newPos]; ok {
				continue
			}
			if _, ok := p.stone[newPos]; ok {
				continue
			}
			// part two sets a floor for us
			// we use the lowest stone Y, although we add the +2 to it in p2
			if p.floor {
				if newPos.Y == p.lowestStoneY {
					continue
				}
			}
			if pos.Y > p.lowestStoneY {
				p.hitAbyss = true
				delete(p.falling, pos)
			} else {
				delete(p.falling, pos)
				p.falling[newPos] = 0
				moved = true
				break
			}
		}
		if !moved {
			delete(p.falling, pos)
			p.atRest[pos] = 0
		}
	}
}

/*
move objects from atRest to atRestBuried if they have objects from atRest or atRestBuried
to their top, top-left and top-right.
*/
func (p *PlayField) Bury() {
	for pos, _ := range p.atRest {
		var count int = 0
		for _, delta := range []Pos{{0, -1}, {-1, -1}, {1, -1}} {
			newPos := pos
			newPos.Add(delta)
			if _, ok := p.stone[newPos]; ok {
				count++
			}
			if _, ok := p.atRest[newPos]; ok {
				count++
			}
			if _, ok := p.atRestBuried[newPos]; ok {
				count++
			}
		}
		if count == 3 {
			delete(p.atRest, pos)
			p.atRestBuried[pos] = 0
		}
	}
}

/*
create a new PlayField.
parse lines from Args[1] using os.ReadFile.
parse the stone lines.
add a new emitter at {500, 0}.
add new falling objects every time there are no falling objects.
move the falling objects until the hitAbyss flag is set.
print out the count of atRest objects as "Part1:"
*/
func main() {
	// create a new PlayField
	var playField PlayField = PlayField{
		emitter:      Pos{500, 0},
		atRest:       make(map[Pos]int),
		atRestBuried: make(map[Pos]int),
		stone:        make(map[Pos]int),
		falling:      make(map[Pos]int),
		lowestStoneY: 0,
		floor:        false,
		hitAbyss:     false,
	}
	// parse lines from Args[1] using os.ReadFile
	lines, err := os.ReadFile(os.Args[1])
	Fatal(err)
	// parse the stone lines
	for _, line := range strings.Split(string(lines), "\n") {
		if line == "" {
			continue
		}
		playField.ParseStoneFromLine(line)
	}
	// add new falling objects every time there are no falling objects
	for !playField.hitAbyss {
		if len(playField.falling) == 0 {
			playField.AddNewFalling()
		}
		playField.MoveFalling()
	}
	// print out the count of atRest objects as "Part1:"
	fmt.Printf("Part1: %d\n", len(playField.atRest))

	// part2
	playField.lowestStoneY += 2
	playField.floor = true
	for {
		if _, ok := playField.atRest[playField.emitter]; ok {
			break
		}
		if len(playField.falling) == 0 {
			playField.Bury()
			playField.AddNewFalling()
		}
		playField.MoveFalling()
	}
	playField.Draw()
	fmt.Printf("Part2: %d\n", len(playField.atRest)+len(playField.atRestBuried))
}
