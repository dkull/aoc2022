package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
)

func Fatal(err error) {
	if err != nil {
		panic(err)
	}
}

type BoundingBox struct {
	x1, x2, y1, y2, z1, z2 int
}

type Pos struct {
	X, Y, Z int
}

type Cube struct {
	pos     Pos
	covered int
	top     bool
	bottom  bool
	left    bool
	right   bool
	front   bool
	back    bool
}

func (c *Cube) CountMarked() int {
	count := 0
	if c.top {
		count++
	}
	if c.bottom {
		count++
	}
	if c.left {
		count++
	}
	if c.right {
		count++
	}
	if c.front {
		count++
	}
	if c.back {
		count++
	}
	return count
}

func (c *Cube) Mark(others map[Pos]*Cube) {
	if _, ok := others[Pos{c.pos.X, c.pos.Y, c.pos.Z + 1}]; ok {
		c.top = true
		others[Pos{c.pos.X, c.pos.Y, c.pos.Z + 1}].bottom = true
	}
	if _, ok := others[Pos{c.pos.X, c.pos.Y, c.pos.Z - 1}]; ok {
		c.bottom = true
		others[Pos{c.pos.X, c.pos.Y, c.pos.Z - 1}].top = true
	}
	if _, ok := others[Pos{c.pos.X, c.pos.Y + 1, c.pos.Z}]; ok {
		c.back = true
		others[Pos{c.pos.X, c.pos.Y + 1, c.pos.Z}].front = true
	}
	if _, ok := others[Pos{c.pos.X, c.pos.Y - 1, c.pos.Z}]; ok {
		c.front = true
		others[Pos{c.pos.X, c.pos.Y - 1, c.pos.Z}].back = true
	}
	if _, ok := others[Pos{c.pos.X + 1, c.pos.Y, c.pos.Z}]; ok {
		c.right = true
		others[Pos{c.pos.X + 1, c.pos.Y, c.pos.Z}].left = true
	}
	if _, ok := others[Pos{c.pos.X - 1, c.pos.Y, c.pos.Z}]; ok {
		c.left = true
		others[Pos{c.pos.X - 1, c.pos.Y, c.pos.Z}].right = true
	}
}

func MarkNeighboringCubes(cubes map[Pos]*Cube, pos Pos) {
	fakeCube := Cube{pos, 0, false, false, false, false, false, false}
	fakeCube.Mark(cubes)
}

func determineBoundingBox(cubes map[Pos]*Cube) BoundingBox {
	// FIXME: This padding is a hack, it won't work without enough padding
	// don't know why atm
	const padding = 10
	x1, y1, z1 := math.MaxInt32, math.MaxInt32, math.MaxInt32
	x2, y2, z2 := math.MinInt32, math.MinInt32, math.MinInt32
	for pos := range cubes {
		if pos.X < x1 {
			x1 = pos.X - padding
		}
		if pos.X > x2 {
			x2 = pos.X + padding
		}
		if pos.Y < y1 {
			y1 = pos.Y - padding
		}
		if pos.Y > y2 {
			y2 = pos.Y + padding
		}
		if pos.Z < z1 {
			z1 = pos.Z - padding
		}
		if pos.Z > z2 {
			z2 = pos.Z + padding
		}
	}
	return BoundingBox{x1, x2, y1, y2, z1, z2}
}

/*
2,2,6
1,2,5
3,2,5
2,1,5
2,3,5
*/
func CubeFromLine(line string) Cube {
	var x, y, z int
	// parse integers from string
	fmt.Sscanf(line, "%d,%d,%d\n", &x, &y, &z)
	return Cube{Pos{x, y, z}, 0, false, false, false, false, false, false}
}

func propagateSteam(bb BoundingBox, cubes map[Pos]*Cube, steam *map[Pos]bool) {
	for {
		addedNewSteam := false
		for pos, exhausted := range *steam {
			if exhausted {
				continue
			}
			MarkNeighboringCubes(cubes, pos)
			(*steam)[pos] = true

			neighborCoords := []Pos{
				{pos.X, pos.Y, pos.Z + 1},
				{pos.X, pos.Y, pos.Z - 1},
				{pos.X, pos.Y + 1, pos.Z},
				{pos.X, pos.Y - 1, pos.Z},
				{pos.X + 1, pos.Y, pos.Z},
				{pos.X - 1, pos.Y, pos.Z},
			}
			for _, neighbor := range neighborCoords {
				// check that neighbor is in bounding box
				if neighbor.X < bb.x1 || neighbor.X > bb.x2 || neighbor.Y < bb.y1 || neighbor.Y > bb.y2 || neighbor.Z < bb.z1 || neighbor.Z > bb.z2 {
					continue
				}
				// check that neighbor is not already in steam
				if _, ok := (*steam)[neighbor]; ok {
					continue
				}
				// check that neighbor is not a cube
				if _, ok := cubes[neighbor]; ok {
					continue
				}
				// add neighbor to new steam
				(*steam)[neighbor] = false
				addedNewSteam = true
			}
		}
		if !addedNewSteam {
			break
		}
	}
}

func main() {
	// read file Argv[1] using os.ReadFile,
	data, err := os.ReadFile(os.Args[1])
	lines := bytes.Split(data, []byte{'\n'})
	Fatal(err)
	// remove last empty line
	lines = lines[:len(lines)-1]

	// parse the lines of form "<int>,<int>,<int> into a a map of Cubes with their coordinates as keys
	cubes := make(map[Pos]*Cube)
	for _, line := range lines {
		// parse the line into a Cube
		cube := CubeFromLine(string(line))
		// add the cube to the map
		cubes[cube.pos] = &cube
	}

	for _, cube := range cubes {
		cube.Mark(cubes)
	}

	freeSide := 0
	for _, cube := range cubes {
		freeSide += 6 - cube.CountMarked()
	}
	fmt.Println("part1:", freeSide)

	// part2

	cubes2 := make(map[Pos]*Cube)
	for _, line := range lines {
		// parse the line into a Cube
		cube := CubeFromLine(string(line))
		// add the cube to the map
		cubes2[cube.pos] = &cube
	}
	bbox := determineBoundingBox(cubes2)
	steamDroplets := make(map[Pos]bool)
	steamDroplets[Pos{bbox.x1, bbox.y1, bbox.z1}] = false
	propagateSteam(bbox, cubes2, &steamDroplets)
	freeSide2 := 0
	for _, cube := range cubes2 {
		freeSide2 += cube.CountMarked()
	}
	fmt.Println("part2:", freeSide2)
}
