package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

/*
Utils
*/

func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

/*
Datastructures
*/

type Pos struct {
	X int
	Y int
}

func (p *Pos) ManhattanDistance(other Pos) int {
	return Abs(p.X-other.X) + Abs(p.Y-other.Y)
}

type Vector struct {
	X int
	Y int
}

type Rule struct {
	FromCubeIdx int
	EntryFacing int
	ToCubeIdx   int
	ExitFacing  int
}

type Player struct {
	PositionX int
	PositionY int
	Facing    int
	MovesLeft int
	Rules     string
	Map       [][]byte
}

type Teleporter struct {
	Position    Pos
	ExitFacing  []int // corner nodes share node, so we need a shared teleporter for them
	TeleportsTo *Teleporter
}

/*
print the map with the player on it.
take additional arguments to only print areaSize*areaSize area around the player.
draw the player as 'P' and the current facing as '^', 'v', '<' or '>'
*/
func (p *Player) PrintMap(areaSize *int) {
	if areaSize == nil {
		areaSize = new(int)
		*areaSize = 10
	}
	for yIdx, line := range p.Map {
		if yIdx < p.PositionY-*areaSize || yIdx > p.PositionY+*areaSize {
			continue
		}
		for xIdx, char := range line {
			if xIdx < p.PositionX-*areaSize || xIdx > p.PositionX+*areaSize {
				continue
			}
			if xIdx == p.PositionX && yIdx == p.PositionY {
				switch p.Facing {
				case 0:
					fmt.Print(">")
				case 1:
					fmt.Print("v")
				case 2:
					fmt.Print("<")
				case 3:
					fmt.Print("^")
				}
			} else {
				fmt.Print(string(char))
			}
		}
		fmt.Println()
	}
}

func (p *Player) ApplyMovementVector(movement Vector) {
	p.PositionX += movement.X
	p.PositionY += movement.Y
}

/*
starting position is the first '#' in the first line
*/
func (p *Player) MoveToStartingPosition() {
	for yIdx, char := range p.Map {
		for xIdx, char := range char {
			if char == '.' {
				p.PositionX = xIdx
				p.PositionY = yIdx
				return
			}
		}
	}
}

/*
PopRule returns the first rule from the rules string. It needs to determine
if it's [\d]+ or [L|R]{1}. Remove that part from the beginning and return it.
*/
func (p *Player) PopRule() (rule string) {
	if len(p.Rules) == 0 {
		fmt.Println("No more rules left")
		return
	}
	if p.Rules[0] >= '0' && p.Rules[0] <= '9' {
		for i, char := range p.Rules {
			if char >= '0' && char <= '9' {
				rule += string(char)
			} else {
				p.Rules = p.Rules[i:]
				break
			}
			if i == len(p.Rules)-1 {
				p.Rules = ""
			}
		}
	} else {
		rule = string(p.Rules[0])
		p.Rules = p.Rules[1:]
	}
	return rule
}

/*
Calculate next tile uses the map and the current position and facing.
It returns the next tile and the next facing.

/*
DoMoves moves the player according to the rules.
This is in a for loop that goes on until there are no more rules left.
If we have a heading and moves left, we move the player until we are out of moves.
Then we pop a new rule, set the heading in relation to our current heading and set the moves left.
If we hit a '#' we stop moving and pop a new rule.
*/
func (p *Player) DoMoves(cubeRules *[]Rule) {
	for {
		if p.MovesLeft == 0 {
			rule := p.PopRule()
			//fmt.Println("popped rule:", rule)
			if rule == "" {
				return
			}
			// if rule is digit, set moves left
			if rule[0] >= '0' && rule[0] <= '9' {
				var err error
				p.MovesLeft, err = strconv.Atoi(rule)
				Fatal(err)
			} else {
				p.Facing = (p.Facing + 1) % 4
				if rule == "L" {
					p.Facing = (p.Facing + 2) % 4
				}
			}
		}
	forLoop1:
		for {
			if p.MovesLeft == 0 {
				break
			}
			// we might need to revert
			savedPositionX := p.PositionX
			savedPositionY := p.PositionY

			//fmt.Println("MOVE FROM p.PositionX:", p.PositionX, "p.PositionY:", p.PositionY, "p.Facing:", p.Facing, "p.MovesLeft:", p.MovesLeft)
			var vector Vector
			switch p.Facing {
			case 0:
				vector = Vector{1, 0}
			case 1:
				vector = Vector{0, 1}
			case 2:
				vector = Vector{-1, 0}
			case 3:
				vector = Vector{0, -1}
			// handle invalid facing
			default:
				panic("invalid facing")
			}

			// move the player
			p.ApplyMovementVector(vector)

			// check where we ended up
			currentTile := p.Map[p.PositionY][p.PositionX]
			//fmt.Println("MOVED TO currentTile:", string(currentTile))
			switch currentTile {
			case '.':
				// good move, just decrement moves left
				p.MovesLeft--
				break forLoop1
			case '#':
				// bad move, revert and pop a new rule
				p.PositionX = savedPositionX
				p.PositionY = savedPositionY
				p.MovesLeft = 0
				//fmt.Println("REVERTING TO p.PositionX:", p.PositionX, "p.PositionY:", p.PositionY, "p.Facing:", p.Facing, "p.MovesLeft:", p.MovesLeft)
				break forLoop1
			case ' ':
				// part1
				if cubeRules == nil {
					//fmt.Println("WRAPPING AROUND")
					// wrap around the map, unless we wrap around to '#'
					// we need to find the next non-' ' tile to determine the success of the move
					wrappingAround := currentTile == ' '
					if wrappingAround {
						// move to opposite side of the map
						switch p.Facing {
						case 0:
							p.PositionX = 0
						case 1:
							p.PositionY = 0
						case 2:
							p.PositionX = len(p.Map[0]) - 1
						case 3:
							p.PositionY = len(p.Map) - 1
						}

						for {
							// check where we ended up
							currentTile = p.Map[p.PositionY][p.PositionX]
							//fmt.Println("in small loop:", "p.PositionX:", p.PositionX, "p.PositionY:", p.PositionY, "p.Facing:", p.Facing, "p.MovesLeft:", p.MovesLeft, "currentTile:", string(currentTile), currentTile)
							switch currentTile {
							case '.':
								//fmt.Println("WRAPPED AROUND TO '.'")
								// good move, just decrement moves left
								p.MovesLeft--
								break forLoop1
							case '#':
								// bad move, revert and pop a new rule
								p.PositionX = savedPositionX
								p.PositionY = savedPositionY
								p.MovesLeft = 0
								break forLoop1
							case ' ':
								p.ApplyMovementVector(vector)
								// keep going
								continue
							}
						}
					}
				} else {
					// part2

					// part 2 requires us to stand on originating block
					switch p.Facing {
					case 0:
						p.PositionX -= 1
					case 1:
						p.PositionY -= 1
					case 2:
						p.PositionX += 1
					case 3:
						p.PositionY += 1
					}
					posX := p.PositionX - 1 // simplify our padding
					posY := p.PositionY - 1 // simplify our padding

					//time.Sleep(1 * time.Second)
					currentTile = p.Map[p.PositionY][p.PositionX]
					fmt.Println("BEFORE TELEPORT")
					fmt.Println("MOVING FROM p.PositionX:", posX, "p.PositionY:", posY, "p.Facing:", p.Facing, "p.MovesLeft:", p.MovesLeft, currentTile)
					areaSize := 30
					p.PrintMap(&areaSize)

					myOrigCubeX := (savedPositionX - 1) / 50
					myOrigCubeY := (savedPositionY - 1) / 50
					origCubeIdx := (myOrigCubeY * 3) + myOrigCubeX // 3 squares in a row
					fmt.Println("Part2 Teleport! from cube:", origCubeIdx)
					newExitFacing := 0
					for ruleIdx := range *cubeRules {
						rule := &(*cubeRules)[ruleIdx]
						if rule.FromCubeIdx != origCubeIdx {
							continue
						}
						if p.Facing != rule.EntryFacing {
							continue
						}
						fmt.Println("rule of interest:", rule)
						exitTopX := (rule.ToCubeIdx % 3 * 50)
						exitTopY := (rule.ToCubeIdx / 3 * 50)
						fmt.Println("from cube Idx:", rule.FromCubeIdx, "to cube Idx:", rule.ToCubeIdx)
						fmt.Println("exitTopX:", exitTopX, "exitTopY:", exitTopY)
						switch rule.EntryFacing { // enter right
						case 0:
							switch rule.ExitFacing {
							case 2: // exit left (WORKS)
								posX = exitTopX + 49
								posY = exitTopY + (49 - ((posY) % 50))
							case 3: // exit up (WORKS)
								posX = exitTopX + ((posY) % 50)
								posY = exitTopY + 49
							default:
								panic("unknown exit facing")
							}
						case 1: // enter down
							switch rule.ExitFacing {
							case 1: // exit down (WORKS)
								posX = exitTopX + (posX % 50)
								posY = exitTopY
							case 2: // exit left (WORKS)
								posY = exitTopY + (posX % 50)
								posX = exitTopX + 49
							default:
								panic("unknown exit facing")
							}
						case 2: // enter left
							switch rule.ExitFacing {
							case 0: // exit right (WORKS)
								posX = exitTopX
								posY = exitTopY + 49 - ((posY) % 50)
							case 1: // exit down (WORDS)
								posX = exitTopX + (posY % 50)
								posY = exitTopY
							default:
								panic("unknown exit facing")
							}
						case 3: // enter up
							switch rule.ExitFacing {
							case 0: // exit right
								posY = exitTopY + (posX % 50)
								posX = exitTopX
							case 3: // exit up
								posX = exitTopX + (posX % 50)
								posY = exitTopY + 49
							default:
								panic("unknown exit facing")
							}
						}
						newExitFacing = rule.ExitFacing
						break
					}

					areaSize = 15

					p.PositionX = posX + 1
					p.PositionY = posY + 1

					currentTile = p.Map[p.PositionY][p.PositionX]
					if currentTile == '#' {
						p.MovesLeft = 0
						p.PositionX = savedPositionX
						p.PositionY = savedPositionY
						break forLoop1
					}
					fmt.Println("Setting new facing to:", newExitFacing, "from:", p.Facing)
					p.Facing = newExitFacing
					p.MovesLeft--
					p.PrintMap(&areaSize)

					if currentTile != '.' {
						fmt.Println("currentTile:", currentTile, "at", posX, posY)
						panic("jump to emptiness")
					}
					// use striing formatting
					switch p.Facing {
					case 0:
						tileToCheck := p.Map[p.PositionY][p.PositionX-1]
						if tileToCheck != ' ' {
							badLandingMsg := fmt.Sprintf("bad landing on %c", tileToCheck)
							panic(badLandingMsg)
						}
						// cases 1-3
					case 1:
						tileToCheck := p.Map[p.PositionY-1][p.PositionX]
						if tileToCheck != ' ' {
							badLandingMsg := fmt.Sprintf("bad landing on %c", tileToCheck)
							panic(badLandingMsg)
						}
					case 2:
						tileToCheck := p.Map[p.PositionY][p.PositionX+1]
						if tileToCheck != ' ' {
							badLandingMsg := fmt.Sprintf("bad landing on %c", tileToCheck)
							panic(badLandingMsg)
						}
					case 3:
						tileToCheck := p.Map[p.PositionY+1][p.PositionX]
						if tileToCheck != ' ' {
							badLandingMsg := fmt.Sprintf("bad landing on %c", tileToCheck)
							panic(badLandingMsg)
						}
					}
					fmt.Println("\n\n\n=============================\n\n\n")
				}
			}
		}
	}
}

/*
Score is 1000*row + 4*col + facing.
We have padded the map with one ' ' on all sides. This means we need to subtract 1 from the row and col.
*/
func (p *Player) GetScore() int {
	// the result coords need to start from 1,1 anyway, so don't subtrack anything
	return 1000*(p.PositionY) + 4*(p.PositionX) + p.Facing
}

/*
Functions
*/

/*
parse a multiline string into lines.
each line consists of ' ', '.' and '#'.
lines may be of different lengths, but they should
be aligned as they are in the input file.
pad shorter lines with ' '.
NOTE!: We pad the map all around with ' ' to make it easier to handle.
*/
func ParseData(input string) (area [][]byte, rules string) {
	twoParts := strings.Split(input, "\n\n")
	mapLines := strings.Split(twoParts[0], "\n")
	rules = twoParts[1]
	// strip ending newline from rules
	rules = rules[:len(rules)-1]

	area = make([][]byte, 0)
	for _, line := range mapLines {
		area = append(area, []byte(line))
	}
	// make all lines the same length by padding shorter lines with ' '
	maxLength := 0
	for _, line := range area {
		if len(line) > maxLength {
			maxLength = len(line)
		}
	}
	for i := range area {
		for len(area[i]) < maxLength {
			area[i] = append(area[i], ' ')
		}
	}
	// pad the map left-right with ' '
	for i := 0; i < len(area); i++ {
		area[i] = append([]byte{' '}, area[i]...)
		area[i] = append(area[i], ' ')
	}
	// add a line of ' ' on top and bottom
	area = append([][]byte{make([]byte, maxLength+2)}, area...)
	area = append(area, make([]byte, maxLength+2))
	// fill them both with ' '
	for i := range area[0] {
		area[0][i] = ' '
		area[len(area)-1][i] = ' '
	}

	return area, rules
}

/*
Making it solve the cube is too much work, i'll just create the teleporters
manually. :'(
*/
func HackyCheat() []Rule {
	rules := []Rule{
		{1, 3, 9, 0}, // green 1
		{9, 2, 1, 1}, // green 2
		{2, 3, 9, 3}, // pink 1
		{9, 1, 2, 1}, // pink 2 +
		{1, 2, 6, 0}, // yellow1
		{6, 2, 1, 0}, // yellow 2
		{2, 0, 7, 2}, // purple1
		{7, 0, 2, 2}, // purple 2 +
		{4, 2, 6, 1}, // teal 1
		{6, 3, 4, 0}, // teal 2
		{7, 1, 9, 2}, // orange 1
		{9, 0, 7, 3}, // orange 2
		{2, 1, 4, 2}, // blue1
		{4, 0, 2, 3}, // blue2
	}
	return rules
}

/*
Main
102221
*/
func main() {
	// parse input file from Argv[1] using os.ReadFile
	data, err := os.ReadFile(os.Args[1])
	Fatal(err)
	// parse data into area and rules
	area, rules := ParseData(string(data))
	// create player
	player := Player{
		PositionX: 0, // needs to be found
		PositionY: 0, // needs to be found
		Facing:    0, // heading 0 is right
		MovesLeft: 0,
		Rules:     rules,
		Map:       area,
	}
	// find player starting position
	player.MoveToStartingPosition()
	fmt.Println("START player.PositionX:", player.PositionX, "player.PositionY:", player.PositionY, "player.Facing:", player.Facing, "player.MovesLeft:", player.MovesLeft)

	player.DoMoves(nil)
	fmt.Println("FINAL OUTCOME")
	player.PrintMap(nil)
	score := player.GetScore()
	fmt.Println("Part1:", score)

	// Part 2
	player = Player{
		PositionX: 0, // needs to be found
		PositionY: 0, // needs to be found
		Facing:    0, // heading 0 is right
		MovesLeft: 0,
		Rules:     rules,
		Map:       area,
	}
	// find player starting position
	player.MoveToStartingPosition()
	hackyRules := HackyCheat()
	player.DoMoves(&hackyRules)
	score2 := player.GetScore()
	fmt.Println("Part2:", score2)
}
