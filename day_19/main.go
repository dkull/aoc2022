package main

import (
	"fmt"
	"os"
	"strings"
)

type Pair[T any, U any] struct {
	first  T
	second U
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Fatal(err error) {
	if err != nil {
		panic(err)
	}
}

type OreRobot struct {
	oreCost int
}

type ClayRobot struct {
	oreCost int
}

type ObsidianRobot struct {
	oreCost  int
	clayCost int
}

type GeodeRobot struct {
	oreCost      int
	obsidianCost int
}

type Recipe struct {
	Id            int
	OreRobot      OreRobot
	ClayRobot     ClayRobot
	ObsidianRobot ObsidianRobot
	GeodeRobot    GeodeRobot
}

type GameState struct {
	ore                        int
	clay                       int
	obsidian                   int
	geode                      int
	oreRobots                  int
	clayRobots                 int
	obsidianRobots             int
	geodeRobots                int
	oreRobotsInProduction      int
	clayRobotsInProduction     int
	obsidianRobotsInProduction int
	geodeRobotsInProduction    int
}

func (gs GameState) BuyOreRobot(recipe Recipe) *GameState {
	// check if we have enough ore
	if gs.ore >= recipe.OreRobot.oreCost {
		return &GameState{
			ore:                        gs.ore - recipe.OreRobot.oreCost,
			clay:                       gs.clay,
			obsidian:                   gs.obsidian,
			geode:                      gs.geode,
			oreRobots:                  gs.oreRobots,
			clayRobots:                 gs.clayRobots,
			obsidianRobots:             gs.obsidianRobots,
			geodeRobots:                gs.geodeRobots,
			oreRobotsInProduction:      gs.oreRobotsInProduction + 1,
			clayRobotsInProduction:     gs.clayRobotsInProduction,
			obsidianRobotsInProduction: gs.obsidianRobotsInProduction,
			geodeRobotsInProduction:    gs.geodeRobotsInProduction,
		}
	}
	return nil
}

func (gs GameState) BuyClayRobot(recipe Recipe) *GameState {
	// check if we have enough ore
	if gs.ore >= recipe.ClayRobot.oreCost {
		return &GameState{
			ore:                        gs.ore - recipe.ClayRobot.oreCost,
			clay:                       gs.clay,
			obsidian:                   gs.obsidian,
			geode:                      gs.geode,
			oreRobots:                  gs.oreRobots,
			clayRobots:                 gs.clayRobots,
			obsidianRobots:             gs.obsidianRobots,
			geodeRobots:                gs.geodeRobots,
			oreRobotsInProduction:      gs.oreRobotsInProduction,
			clayRobotsInProduction:     gs.clayRobotsInProduction + 1,
			obsidianRobotsInProduction: gs.obsidianRobotsInProduction,
			geodeRobotsInProduction:    gs.geodeRobotsInProduction,
		}
	}
	return nil
}

func (gs GameState) BuyObsidianRobot(recipe Recipe) *GameState {
	// check if we have enough ore
	if gs.ore >= recipe.ObsidianRobot.oreCost && gs.clay >= recipe.ObsidianRobot.clayCost {
		return &GameState{
			ore:                        gs.ore - recipe.ObsidianRobot.oreCost,
			clay:                       gs.clay - recipe.ObsidianRobot.clayCost,
			obsidian:                   gs.obsidian,
			geode:                      gs.geode,
			oreRobots:                  gs.oreRobots,
			clayRobots:                 gs.clayRobots,
			obsidianRobots:             gs.obsidianRobots,
			geodeRobots:                gs.geodeRobots,
			oreRobotsInProduction:      gs.oreRobotsInProduction,
			clayRobotsInProduction:     gs.clayRobotsInProduction,
			obsidianRobotsInProduction: gs.obsidianRobotsInProduction + 1,
			geodeRobotsInProduction:    gs.geodeRobotsInProduction,
		}
	}
	return nil
}

func (gs GameState) BuyGeodeRobot(recipe Recipe) *GameState {
	// check if we have enough ore
	if gs.ore >= recipe.GeodeRobot.oreCost && gs.obsidian >= recipe.GeodeRobot.obsidianCost {
		return &GameState{
			ore:                        gs.ore - recipe.GeodeRobot.oreCost,
			clay:                       gs.clay,
			obsidian:                   gs.obsidian - recipe.GeodeRobot.obsidianCost,
			geode:                      gs.geode,
			oreRobots:                  gs.oreRobots,
			clayRobots:                 gs.clayRobots,
			obsidianRobots:             gs.obsidianRobots,
			geodeRobots:                gs.geodeRobots,
			oreRobotsInProduction:      gs.oreRobotsInProduction,
			clayRobotsInProduction:     gs.clayRobotsInProduction,
			obsidianRobotsInProduction: gs.obsidianRobotsInProduction,
			geodeRobotsInProduction:    gs.geodeRobotsInProduction + 1,
		}
	}
	return nil
}

/*
Parse recipes from a string:
Blueprint 1: Each ore robot costs 4 ore. Each clay robot costs 4 ore. Each obsidian robot costs 4 ore and 8 clay. Each geode robot costs 2 ore and 15 obsidian.
Blueprint 2: Each ore robot costs 4 ore. Each clay robot costs 4 ore. Each obsidian robot costs 3 ore and 19 clay. Each geode robot costs 4 ore and 15 obsidian.
Blueprint 3: Each ore robot costs 4 ore. Each clay robot costs 4 ore. Each obsidian robot costs 2 ore and 8 clay. Each geode robot costs 3 ore and 9 obsidian.
*/
func ParseRecipes(input []string) []Recipe {
	recipes := make([]Recipe, 0)
	for i, line := range input {
		if line == "" {
			continue
		}
		oreRobot := OreRobot{}
		clayRobot := ClayRobot{}
		obsidianRobot := ObsidianRobot{}
		geodeRobot := GeodeRobot{}
		fmt.Println("line:", line)
		_, err := fmt.Sscanf(line, "Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian.", &i, &oreRobot.oreCost, &clayRobot.oreCost, &obsidianRobot.oreCost, &obsidianRobot.clayCost, &geodeRobot.oreCost, &geodeRobot.obsidianCost)
		Fatal(err)
		recipes = append(recipes, Recipe{
			Id:            i,
			OreRobot:      oreRobot,
			ClayRobot:     clayRobot,
			ObsidianRobot: obsidianRobot,
			GeodeRobot:    geodeRobot,
		})
	}
	return recipes
}

var highestGeodes = 0

func simulate(recipe Recipe, gs GameState, minute int, maxminute int) (totalGeodes int) {
	if minute > maxminute {
		return gs.geode
	}

	if gs.geode > highestGeodes {
		highestGeodes = gs.geode
		//fmt.Println("minute:", minute, "geodes:", gs.geode, "highestGeodes:", highestGeodes)
	} else {
		// how many we could magically make
		bestCaseGeodes := gs.geode
		bestCaseMiners := gs.geodeRobots
		for i := minute + 1; i <= maxminute; i++ {
			bestCaseGeodes += bestCaseMiners
			bestCaseMiners += 1
		}
		if bestCaseGeodes < highestGeodes {
			return 0
		}
		// do i have obsidian to build a geode robot?
		/*
			if gs.geode+(gs.geodeRobots*maxminute-minute) < highestGeodes {
				// no, we can't win with our current geode robots
				// check if we build a new obsidian robot every turn, could we win?
				obsidian := gs.obsidian
				obsidianRobots := gs.obsidianRobots
				for i := minute + 1; i <= maxminute; i++ {
					obsidian += obsidianRobots
					obsidianRobots += 1
				}
				if obsidian < recipe.GeodeRobot.obsidianCost {
					// no, we can't win
					return 0
				}
			}*/
	}

	// if any robots are in production, we skip the production step
	notProducing := !(gs.oreRobotsInProduction > 0 || gs.clayRobotsInProduction > 0 || gs.obsidianRobotsInProduction > 0 || gs.geodeRobotsInProduction > 0)

	bestBranchResult := 0
	if notProducing && maxminute-minute >= 1 {
		if newState := gs.BuyGeodeRobot(recipe); newState != nil {
			if geodeCount := simulate(recipe, *newState, minute, maxminute); geodeCount > bestBranchResult {
				bestBranchResult = geodeCount
			}
		}
		if newState := gs.BuyObsidianRobot(recipe); newState != nil {
			if geodeCount := simulate(recipe, *newState, minute, maxminute); geodeCount > bestBranchResult {
				bestBranchResult = geodeCount
			}
		}
		if newState := gs.BuyClayRobot(recipe); newState != nil {
			if geodeCount := simulate(recipe, *newState, minute, maxminute); geodeCount > bestBranchResult {
				bestBranchResult = geodeCount
			}
		}
		if newState := gs.BuyOreRobot(recipe); newState != nil {
			if geodeCount := simulate(recipe, *newState, minute, maxminute); geodeCount > bestBranchResult {
				bestBranchResult = geodeCount
			}
		}
	}
	// collect stuff
	gs.ore += gs.oreRobots
	gs.clay += gs.clayRobots
	gs.obsidian += gs.obsidianRobots
	gs.geode += gs.geodeRobots
	// produce robots
	gs.oreRobots += gs.oreRobotsInProduction
	gs.oreRobotsInProduction = 0
	gs.clayRobots += gs.clayRobotsInProduction
	gs.clayRobotsInProduction = 0
	gs.obsidianRobots += gs.obsidianRobotsInProduction
	gs.obsidianRobotsInProduction = 0
	gs.geodeRobots += gs.geodeRobotsInProduction
	gs.geodeRobotsInProduction = 0
	simulationResult := simulate(recipe, gs, minute+1, maxminute)
	//fmt.Printf("%d[%d]: %+v\n", minute, simulationResult, gs)
	return Max(bestBranchResult, simulationResult)
}

/*
open file Argv[1] using os.ReadFile.
split data into lines, remove last empty line.
parse recipes from lines.
*/
func main() {
	data, err := os.ReadFile(os.Args[1])
	Fatal(err)
	lines := strings.Split(string(data), "\n")
	recipes := ParseRecipes(lines)
	fmt.Println(recipes)

	qualityLvlSum := 0
	for _, recipe := range recipes {
		highestGeodes = 0
		gamestate := GameState{
			ore:            0,
			clay:           0,
			obsidian:       0,
			geode:          0,
			oreRobots:      1,
			clayRobots:     0,
			obsidianRobots: 0,
			geodeRobots:    0,
		}
		result := simulate(recipe, gamestate, 1, 24)
		qualityLvl := result * recipe.Id
		qualityLvlSum += qualityLvl
		fmt.Println("blueprint", recipe.Id, "result:", result, "new quality lvl:", qualityLvl, "total quality lvl:", qualityLvlSum)
	}

	for _, recipe := range recipes {
		highestGeodes = 1
		gamestate := GameState{
			ore:            0,
			clay:           0,
			obsidian:       0,
			geode:          0,
			oreRobots:      1,
			clayRobots:     0,
			obsidianRobots: 0,
			geodeRobots:    0,
		}
		result := simulate(recipe, gamestate, 1, 32)
		highestGeodes *= result
		fmt.Println("blueprint", recipe.Id, "result:", result)
	}
}
