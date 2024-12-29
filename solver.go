package main

import (
	"fmt"
	"strings"
)

type SearchType int

const (
	DFSSearchType SearchType = iota
	BFSSearchType
)

type Solver struct{}

func NewSolver() *Solver {
	return &Solver{}
}

func (s *Solver) SolveMaze(m *MazeGenerator, mazeStr string, searchType SearchType) string {
	if mazeStr != "" {
		m.originalMaze = mazeStr
	}
	m.clearMaze()

	maze, start, end, err := s.parseMaze(m.originalMaze)

	if err != nil {
		return err.Error()
	}

	visited := make([][]bool, len(maze))
	parent := make(map[Point]Point) // for BFS
	for i := range visited {
		visited[i] = make([]bool, len(maze[0])/2)
	}

	if m.fastMode {
		m.isSolving = false
		path := make(map[Point]bool)
		var success bool
		if searchType == DFSSearchType {
			success = s.dfs(m, maze, start, end, visited, path, nil, nil)
		} else {
			success = s.bfs(m, maze, start, end, visited, path, parent, nil, nil)
		}
		if success {
			return s.visualizeExploration(m, maze, path)
		}
		return "No solution found"
	}

	path := make(map[Point]bool)
	updates := make(chan string)
	var finalSolution string

	// Printing out the steps and log in seperate of the main thread
	go func() {
		defer close(updates)
		defer func() {
			m.isSolving = false
		}()

		stepCount := 0

		var success bool
		if searchType == DFSSearchType {
			m.logger.LogStart("DFS")
			success = s.dfs(m, maze, start, end, visited, path, updates, func() {
				stepCount++
				m.logger.LogStep(stepCount, start.col, start.row, false)
			})
			m.logger.LogComplete(stepCount, 0)
		} else {
			m.logger.LogStart("BFS")
			success = s.bfs(m, maze, start, end, visited, path, parent, updates, func() {
				stepCount++
				m.logger.LogStep(stepCount, start.col, start.row, false)
			})
			m.logger.LogComplete(stepCount, 0)
		}

		if !success {
			m.logger.LogNoSolution()
		}
	}()

	for update := range updates {
		finalSolution = update
	}

	return finalSolution
}

func (s *Solver) dfs(m *MazeGenerator, maze [][]rune, current, end Point, visited [][]bool,
	path map[Point]bool, updates chan<- string, onStep func()) bool {

	if onStep != nil {
		onStep()
	}

	if current.row < 0 || current.row >= len(maze) ||
		current.col*2 >= len(maze[0]) ||
		maze[current.row][current.col*2] == '█' ||
		visited[current.row][current.col] {
		return false
	}

	visited[current.row][current.col] = true
	path[current] = true

	if updates != nil {
		updates <- s.visualizeExploration(m, maze, path)
	}

	if current == end {
		return true
	}

	directions := []Point{
		{0, -1}, // up
		{1, 0},  // right
		{0, 1},  // down
		{-1, 0}, // left
	}

	for _, dir := range directions {
		next := Point{current.col + dir.col, current.row + dir.row}
		if s.dfs(m, maze, next, end, visited, path, updates, onStep) {
			return true
		}
	}

	delete(path, current)

	if updates != nil {
		updates <- s.visualizeExploration(m, maze, path)
	}
	return false
}

func (s *Solver) bfs(m *MazeGenerator, maze [][]rune, start, end Point, visited [][]bool,
	path map[Point]bool, parent map[Point]Point, updates chan<- string, onStep func()) bool {

	queue := []Point{start}
	visited[start.row][start.col] = true
	explored := make(map[Point]bool)

	directions := []Point{
		{0, -1}, // up
		{1, 0},  // right
		{0, 1},  // down
		{-1, 0}, // left
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if onStep != nil {
			onStep()
		}

		explored[current] = true

		if updates != nil {
			updates <- s.visualizeExploration(m, maze, explored)
		}

		if current == end {
			for k := range path {
				delete(path, k)
			}
			finalPath := reconstructPath(parent, end)
			for p := range finalPath {
				path[p] = true
			}
			if updates != nil {
				updates <- s.visualizeExploration(m, maze, path)
			}
			return true
		}

		for _, dir := range directions {
			next := Point{current.col + dir.col, current.row + dir.row}
			if next.row >= 0 && next.row < len(maze) &&
				next.col >= 0 && next.col < len(maze[0])/2 &&
				!visited[next.row][next.col] &&
				maze[next.row][next.col*2] != '█' {
				visited[next.row][next.col] = true
				parent[next] = current
				queue = append(queue, next)
			}
		}
	}

	return false
}

func (s *Solver) parseMaze(mazeStr string) ([][]rune, Point, Point, error) {
	var maze [][]rune
	lines := strings.Split(mazeStr, "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		maze = append(maze, []rune(line))
	}

	if len(maze) == 0 || len(maze[0]) == 0 {
		return nil, Point{}, Point{}, fmt.Errorf("invalid maze")
	}

	rows, cols := len(maze), len(maze[0])
	var start, end Point
	for col := 0; col < cols; col += 2 {
		if col+1 < cols && maze[0][col] == ' ' && maze[0][col+1] == ' ' {
			start = Point{col / 2, 0}
		}
		if col+1 < cols && maze[rows-1][col] == ' ' && maze[rows-1][col+1] == ' ' {
			end = Point{col / 2, rows - 1}
		}
	}

	return maze, start, end, nil
}

// Visualize the green path we explore to find the solution
func (s *Solver) visualizeExploration(m *MazeGenerator, maze [][]rune, path map[Point]bool) string {
	blocks := BlockStyles[m.blockStyle]
	result := ""

	for i := range maze {
		for j := 0; j < len(maze[i]); j += 2 {
			pos := Point{j / 2, i}
			if path[pos] {
				result += "[green]" + blocks.Solution + "[white]"
			} else if j+1 < len(maze[i]) {
				result += string([]rune{maze[i][j], maze[i][j+1]})
			}
		}
		result += "\n"
	}

	if m.app != nil {
		m.app.QueueUpdateDraw(func() {
			m.matrixBox.SetText(result)
		})
	}
	return result
}

// (BFS) Backtrack to get the solution path
func reconstructPath(parent map[Point]Point, current Point) map[Point]bool {
	path := make(map[Point]bool)
	for {
		path[current] = true
		next, exists := parent[current]
		if !exists {
			break
		}
		current = next
	}
	return path
}

// Clear the maze
func (m *MazeGenerator) clearMaze() {
	if m.originalMaze == "" {
		return
	}

	maze := strings.Split(m.originalMaze, "\n")
	for i, line := range maze {
		runes := []rune(line)
		for j := 0; j < len(runes); j += 2 {
			if runes[j] != '█' {
				runes[j] = ' '
				if j+1 < len(runes) {
					runes[j+1] = ' '
				}
			}
		}
		maze[i] = string(runes)
	}
	m.originalMaze = strings.Join(maze, "\n")
}
