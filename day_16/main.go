package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func Fatal(err error) {
	if err != nil {
		panic(err)
	}
}

type Valve struct {
	name string
	rate int
}

type Link struct {
	a        string
	b        string
	distance int
}

func (l Link) IsBetween(a string, b string) bool {
	return (l.a == a && l.b == b) || (l.a == b && l.b == a)
}

func (l Link) GetOther(a string) string {
	if l.a == a {
		return l.b
	}
	return l.a
}

func ArrContains[T comparable](arr []T, val T) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func containsValve(valves map[string]int, valve string) bool {
	for k := range valves {
		if k == valve {
			return true
		}
	}
	return false
}

func containsValveArr(valves []string, valve string) bool {
	for _, v := range valves {
		if v == valve {
			return true
		}
	}
	return false
}

func containsLink(links []Link, link Link) bool {
	for _, l := range links {
		if l == link {
			return true
		}
	}
	return false
}

/*
parse lines like this:

	Valve AA has flow rate=0; tunnels lead to valves DD, II, BB
	Valve BB has flow rate=13; tunnels lead to valves CC, AA
	Valve CC has flow rate=13; tunnel leads to valve CC

parse them into Valve structs. and links
the name is Valve <name>. rate is rate=<rate>. paths are tunnels lead to valves <name1>, <name2>, ...
*/
func parseValve(line string) (Valve, []Link) {
	// Use a regular expression to extract the name, rate, and paths
	re := regexp.MustCompile(`Valve (\w+) has flow rate=(\d+); tunnel.? lead.? to valve.? (.*)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) < 3 {
		// Return an empty valve if the regular expression didn't match
		panic("Invalid valve line: " + line)
	}

	// Extract the name, rate, and paths from the matches
	name := matches[1]
	rate, _ := strconv.Atoi(matches[2])
	paths := strings.Split(matches[3], ", ")

	// Create links for each path
	links := []Link{}
	for _, path := range paths {
		links = append(links, Link{name, path, 1})
	}

	// Create and return a new Valve struct
	return Valve{
		name: name,
		rate: rate,
	}, links
}

/*
simplify the graph of Valves. remove all valves with rate=0, then increment all the neighboring
valves distance by 1, and link the neighbors together. repeat until no more valves with rate=0.
we will have rate=0 valves as neighbors, so we can't delete more than one valve at a time.
*/
func simplifyGraph(valves map[string]Valve, links []Link) (map[string]Valve, []Link) {
	for {
		// find a valve with rate=0 and two neighbors
		// remove the valve and link the neighbors together
		// repeat until no more valves with rate=0
		removed := false
		for valveName, valve := range valves {
			if valve.rate != 0 || valve.name == "AA" {
				continue
			}
			valveLinks := []Link{}
			for _, link := range links {
				if link.a == valveName || link.b == valveName {
					valveLinks = append(valveLinks, link)
				}
			}
			for _, link1 := range valveLinks {
				for _, link2 := range valveLinks {
					if link1 == link2 {
						continue
					}
					neighbor1 := link1.GetOther(valveName)
					neighbor2 := link2.GetOther(valveName)
					if neighbor1 == neighbor2 {
						continue
					}

					newDistance := link1.distance + link2.distance
					linkExists := false
					for _, existingLink := range links {
						// check if these two already have a link
						if existingLink.IsBetween(neighbor1, neighbor2) {
							linkExists = true
							// if they do, update the distances to be the smallest known distance
							if existingLink.distance > newDistance {
								existingLink.distance = newDistance
							}
							// do not append a new link if we found an existing one
							continue
						}
					}
					if !linkExists {
						links = append(links, Link{neighbor1, neighbor2, link1.distance + link2.distance})
					}

					removed = true
					delete(valves, valveName)
				}
			}
		}
		if !removed {
			break
		}
	}
	return valves, links
}

func pruneDuplicateLinks(valves map[string]Valve, links []Link) []Link {
	// remove duplicate links
	uniqueLinks := []Link{}
	for _, link := range links {
		// if valves does not contain either link endpoint, remove the link
		if _, ok := valves[link.a]; !ok {
			continue
		}
		if _, ok := valves[link.b]; !ok {
			continue
		}
		// see if colliding
		collision := false
		for _, uniqueLink := range uniqueLinks {
			if link.a == uniqueLink.a && link.b == uniqueLink.b {
				collision = true
				break
			}
			if link.a == uniqueLink.b && link.b == uniqueLink.a {
				collision = true
				break
			}
		}
		if !collision {
			uniqueLinks = append(uniqueLinks, link)
		}
	}
	return uniqueLinks
}

/*
use link distance to find the shortest path from AA to ZZ
*/
func findShortestPath(a Valve, b Valve, valves map[string]Valve, visited []string, links []Link) int {
	// if we have visited this valve, return a high number
	if containsValveArr(visited, a.name) {
		return 100000
	}

	// if we have reached the end, return the distance
	if a.name == b.name {
		return 0
	}

	// find all the links for this valve
	valveLinks := []Link{}
	for _, link := range links {
		if link.a == a.name || link.b == a.name {
			valveLinks = append(valveLinks, link)
		}
	}

	// find the shortest path from each neighbor
	visited = append(visited, a.name)
	shortest := 100000
	for _, link := range valveLinks {
		neighbor := link.GetOther(a.name)
		distance := link.distance + findShortestPath(valves[neighbor], b, valves, visited, links)
		if distance < shortest {
			shortest = distance
		}
	}

	return shortest
}

/*
given a current valve, and a target valve, calculate the flow rate we would achieve by the end if
we chose to go to the target valve and turn it. the valves form a graph, so we need to move through
other valves to get to the target valve. each move costs 1 minute, and each valve turning costs 1 minute.
we cannot exceed minutesLeft. we need to use the shortest path when moving to target valve.
return the valves from rate multiplied by the minutes we have left after reaching the target valve.
*/
func calculateFlowRate(valves map[string]Valve, linkmap map[string]map[string]int, opened []string, minutesLeft int, currentValve Valve) (int, []string) {
	// If we have no minutes left, return 0
	if minutesLeft <= 0 {
		return 0, []string{}
	}

	// turn the valve
	minutesLeft -= 1 // turn the valve
	currentValveScore := currentValve.rate * minutesLeft
	opened = append(opened, currentValve.name)

	bestScore := 0
	bestRoute := []string{}
	valvePaths := linkmap[currentValve.name]
	for neighborValveName, distance := range valvePaths {
		// if we have already opened this valve, skip it
		if ArrContains(opened, neighborValveName) {
			continue
		}
		// if we have no minutes left to reach the neighbor and turn it, skip it
		if minutesLeft < distance+1 {
			continue
		}
		score, route := calculateFlowRate(valves, linkmap, opened, minutesLeft-distance, valves[neighborValveName])
		if score > bestScore {
			bestScore = score
			bestRoute = route
		}
	}

	return currentValveScore + bestScore, append([]string{currentValve.name}, bestRoute...)
}

/*
read a file using os.ReadFile from Args[1]. parse the file into Valve structs. for each valve,
calculate the maximum flow rate we can achieve by turning the valve. print the maximum flow rate
*/
func main() {
	// Read the file
	file, err := os.ReadFile(os.Args[1])
	Fatal(err)

	// Split the file into lines
	lines := strings.Split(string(file), "\n")
	// remove the last line if it is empty
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	// Parse the lines into Valve structs
	valves := make(map[string]Valve)
	links := []Link{}
	for _, line := range lines {
		valve, valveLinks := parseValve(line)
		links = append(links, valveLinks...)
		if valve.name != "" {
			valves[valve.name] = valve
		}
	}

	fmt.Println("before simplification", len(valves), "valves", len(links), "links")
	valves, links = simplifyGraph(valves, links)
	fmt.Println("after simplification", len(valves), "valves", len(links), "links")
	// prune links where one is <a,b> and other is <b,a>
	fmt.Println("pruning duplicate links (before pruning)", len(links))
	links = pruneDuplicateLinks(valves, links)
	fmt.Println("pruning duplicate links (after pruning)", len(links))

	// create a distance mapping between all Valves, this gives
	// as the simplest possible way to move from one valve to another
	distanceMap := make(map[string]map[string]int)
	for _, valve := range valves {
		distanceMap[valve.name] = make(map[string]int)
		for _, valve2 := range valves {
			shortestPath := findShortestPath(valve, valve2, valves, []string{}, links)
			distanceMap[valve.name][valve2.name] = shortestPath
		}
	}

	// print distance map
	fmt.Println("distance map")
	for _, valve := range valves {
		fmt.Println(valve.name, distanceMap[valve.name])
	}

	// Calculate the maximum flow rate for each valve from valve 'AA'
	opened := []string{}
	atValve := valves["AA"]
	bestFlowRate, bestRoute := calculateFlowRate(valves, distanceMap, opened, 31, atValve)
	fmt.Println("Part1:", bestFlowRate, bestRoute)

	bestFlowRate, bestRoute = calculateFlowRate(valves, distanceMap, opened, 10, atValve)
	fmt.Println("Part2:", bestFlowRate, bestRoute)
}
