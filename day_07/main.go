package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Directory struct {
	Fullpath      string
	ChildrenNames []string
	Files         []string
	DeepSize      int
	ShallowSize   int
	Listed        bool
}

func ChopDirOff(fullpath string) string {
	if strings.Count(fullpath, "/") == 1 {
		return "/"
	}
	return strings.Join(strings.Split(fullpath, "/")[:len(strings.Split(fullpath, "/"))-1], "/")
}

func PushOnDir(fullpath string, dir string) string {
	if fullpath == "/" {
		return "/" + dir
	}
	if fullpath == "/" && dir == "/" {
		return "/"
	}
	if fullpath == "" && dir == "/" {
		return "/"
	}
	return fullpath + "/" + dir
}

func HandleChangeDirectory(dirs *map[string]*Directory, currentDir *Directory, newDir string) *Directory {
	if newDir == ".." {
		newDir = ChopDirOff(currentDir.Fullpath)
		currentDir = (*dirs)[newDir]
		return currentDir
	}
	newDir = PushOnDir(currentDir.Fullpath, newDir)
	if _, ok := (*dirs)[newDir]; !ok {
		(*dirs)[newDir] = &Directory{
			Fullpath:      newDir,
			ChildrenNames: []string{},
			Files:         []string{},
			ShallowSize:   0,
			DeepSize:      0,
			Listed:        false,
		}
	}
	return (*dirs)[newDir]
}

/*
Processes a list of simulated CLI directory traversal commands and file listings.
Collects the information received into a map of Directory structs.
Commands starting with '$' are either 'cd' or 'ls' commands. Eg. '$ cd foo' or '$ cd ..'
'cd' command changes the current directory.
'ls' command lists the files and directories in the current directory.
Name the command parts[2] as targetDirName
ls output is in the following format:
dir somedirectory
12345 somefile
...
*/
func ProcessCommands(commands []string) map[string]*Directory {
	dirs := make(map[string]*Directory)
	var currentDir *Directory = &Directory{}
	for _, command := range commands {
		parts := strings.Split(command, " ")
		if parts[0] == "$" {
			if parts[1] == "cd" {
				currentDir = HandleChangeDirectory(&dirs, currentDir, parts[2])
				dirs[currentDir.Fullpath] = currentDir
			} else if parts[1] == "ls" {
				currentDir.Listed = true
			} else {
				log.Fatal("Unknown command: ", parts[1])
			}
		} else {
			if parts[0] == "dir" {
				currentDir.ChildrenNames = append(currentDir.ChildrenNames, parts[1])
			} else {
				currentDir.Files = append(currentDir.Files, parts[1])
				size, err := strconv.Atoi(parts[0])
				Fatal(err)
				currentDir.ShallowSize += size
			}
		}
	}
	return dirs
}

func FindTotalSizes(curPath string, dirs map[string]*Directory) {
	dir := dirs[curPath]
	for _, childName := range dir.ChildrenNames {
		childPath := PushOnDir(curPath, childName)
		FindTotalSizes(childPath, dirs)
		dir.DeepSize += dirs[childPath].DeepSize
	}
	dir.DeepSize += dir.ShallowSize
}

func ReadLines(path string) []string {
	file, err := os.Open(path)
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

func SumDirsByTotalSize(dirs map[string]*Directory, totalSmallerThan int) int {
	sum := 0
	for _, dir := range dirs {
		if dir.DeepSize < totalSmallerThan {
			sum += dir.DeepSize
		}
	}
	return sum
}

/*
This function takes a map of Directory structs, how much disk space we have and how much we need to be free.
We then find one directory which if deleted would free up enough space. We choose the smallest possible directory.
We return the total size of that directory.
We start by iterating over the map and storing the smallest directory size that would be enough.
*/
func FindSmallestDir(dirs map[string]*Directory, diskSize int, neededSpace int) int {
	minSize := diskSize
	// total disk space used
	usedSize := dirs["/"].DeepSize
	freeSize := diskSize - usedSize
	for _, dir := range dirs {
		// find the directory which if deleted would free up enough space
		if freeSize+dir.DeepSize >= neededSpace {
			if dir.DeepSize < minSize {
				minSize = dir.DeepSize
			}
		}
	}
	return minSize
}

/*
Read input file into lines from os.Argv[1]
Feed the lines to ProcessCommands
*/
func main() {
	lines := ReadLines(os.Args[1])
	dirs := ProcessCommands(lines)
	FindTotalSizes("/", dirs)
	fmt.Println("Part1:", SumDirsByTotalSize(dirs, 100000))
	fmt.Println("Part2:", FindSmallestDir(dirs, 70000000, 30000000))

}
