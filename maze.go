package main

import (
	"math/rand"

	"github.com/rivo/tview"
)

type AlgorithmState int
type Point struct{ col, row int }
type MazeAlgorithm int

const (
	DFS MazeAlgorithm = iota
	Prims
	Kruskals
)

type MazeGenerator struct {
	blockStyle   BlockStyle
	startPos     Point
	endPos       Point
	maze         [][]bool
	app          *tview.Application
	matrixBox    *tview.TextView
	logView      *tview.TextView
	logger       *LogManager
	originalMaze string
	fastMode     bool
	isSolving    bool
}

func (m *MazeGenerator) GenerateMaze(totalRows, totalCols int, algorithm MazeAlgorithm) string {
	mazeRows := (totalRows/2)*2 + 1
	mazeCols := (totalCols/2)*2 + 1

	// a plate of marble for me to carve my maze
	maze := make([][]bool, mazeRows)
	for i := range maze {
		maze[i] = make([]bool, mazeCols)
		for j := range maze[i] {
			maze[i][j] = true
		}
	}

	switch algorithm {
	case DFS:
		m.generateDFSMaze(maze, mazeRows, mazeCols)
	case Prims:
		m.generatePrimsMaze(maze, mazeRows, mazeCols)
	case Kruskals:
		m.generateKruskalsMaze(maze, mazeRows, mazeCols)
	}

	entranceCol := 1 + 2*rand.Intn((mazeCols-1)/2)
	exitCol := 1 + 2*rand.Intn((mazeCols-1)/2)

	maze[0][entranceCol] = false
	maze[mazeRows-1][exitCol] = false

	m.maze = maze
	m.startPos = Point{entranceCol, 1}
	m.endPos = Point{exitCol, mazeRows - 2}

	return m.DisplayMaze()
}

func (m *MazeGenerator) generateDFSMaze(maze [][]bool, mazeRows, mazeCols int) {
	startCol, startRow := 1, 1
	maze[startRow][startCol] = false
	stack := []Point{{startCol, startRow}}

	dCol := []int{2, 0, -2, 0}
	dRow := []int{0, 2, 0, -2}

	for len(stack) > 0 {
		current := stack[len(stack)-1]

		var unvisited []int
		for dir := 0; dir < 4; dir++ {
			newCol := current.col + dCol[dir]
			newRow := current.row + dRow[dir]

			if newCol > 0 && newCol < mazeCols-1 && newRow > 0 && newRow < mazeRows-1 && maze[newRow][newCol] {
				unvisited = append(unvisited, dir)
			}
		}

		if len(unvisited) == 0 {
			stack = stack[:len(stack)-1]
			continue
		}

		dir := unvisited[rand.Intn(len(unvisited))]
		newCol := current.col + dCol[dir]
		newRow := current.row + dRow[dir]

		maze[newRow][newCol] = false
		maze[current.row+dRow[dir]/2][current.col+dCol[dir]/2] = false

		stack = append(stack, Point{newCol, newRow})
	}
}

func (m *MazeGenerator) generatePrimsMaze(maze [][]bool, mazeRows, mazeCols int) {
	for i := 0; i < mazeRows; i++ {
		for j := 0; j < mazeCols; j++ {
			maze[i][j] = true
		}
	}

	dCol := []int{0, 2, 0, -2}
	dRow := []int{-2, 0, 2, 0}

	start := Point{1, 1}
	maze[start.row][start.col] = false

	walls := make([]Point, 0)
	for i := 0; i < 4; i++ {
		newCol := start.col + dCol[i]/2
		newRow := start.row + dRow[i]/2
		if newCol > 0 && newCol < mazeCols-1 && newRow > 0 && newRow < mazeRows-1 {
			walls = append(walls, Point{newCol, newRow})
		}
	}

	for len(walls) > 0 {
		wallIdx := rand.Intn(len(walls))
		wall := walls[wallIdx]

		walls = append(walls[:wallIdx], walls[wallIdx+1:]...)

		for i := 0; i < 4; i++ {
			cell1X := wall.col + dCol[i]/2
			cell1Y := wall.row + dRow[i]/2
			cell2X := wall.col - dCol[i]/2
			cell2Y := wall.row - dRow[i]/2

			if cell1X > 0 && cell1X < mazeCols-1 && cell1Y > 0 && cell1Y < mazeRows-1 &&
				cell2X > 0 && cell2X < mazeCols-1 && cell2Y > 0 && cell2Y < mazeRows-1 {
				if !maze[cell1Y][cell1X] && maze[cell2Y][cell2X] {
					maze[wall.row][wall.col] = false
					maze[cell2Y][cell2X] = false

					for j := 0; j < 4; j++ {
						newCol := cell2X + dCol[j]/2
						newRow := cell2Y + dRow[j]/2
						if newCol > 0 && newCol < mazeCols-1 && newRow > 0 && newRow < mazeRows-1 && maze[newRow][newCol] {
							walls = append(walls, Point{newCol, newRow})
						}
					}
					break
				} else if maze[cell1Y][cell1X] && !maze[cell2Y][cell2X] {
					maze[wall.row][wall.col] = false
					maze[cell1Y][cell1X] = false

					for j := 0; j < 4; j++ {
						newCol := cell1X + dCol[j]/2
						newRow := cell1Y + dRow[j]/2
						if newCol > 0 && newCol < mazeCols-1 && newRow > 0 && newRow < mazeRows-1 && maze[newRow][newCol] {
							walls = append(walls, Point{newCol, newRow})
						}
					}
					break
				}
			}
		}
	}
}

func (m *MazeGenerator) generateKruskalsMaze(maze [][]bool, mazeRows, mazeCols int) {
	type DisjointSet struct {
		parent []int
		rank   []int
	}

	newDisjointSet := func(size int) *DisjointSet {
		ds := &DisjointSet{
			parent: make([]int, size),
			rank:   make([]int, size),
		}
		for i := range ds.parent {
			ds.parent[i] = i
		}
		return ds
	}

	var findFunc func(ds *DisjointSet, col int) int
	findFunc = func(ds *DisjointSet, col int) int {
		if ds.parent[col] != col {
			ds.parent[col] = findFunc(ds, ds.parent[col])
		}
		return ds.parent[col]
	}

	union := func(ds *DisjointSet, col, row int) {
		pCol, pRow := findFunc(ds, col), findFunc(ds, row)
		if pCol == pRow {
			return
		}
		if ds.rank[pCol] < ds.rank[pRow] {
			ds.parent[pCol] = pRow
		} else if ds.rank[pCol] > ds.rank[pRow] {
			ds.parent[pRow] = pCol
		} else {
			ds.parent[pRow] = pCol
			ds.rank[pCol]++
		}
	}

	for i := range maze {
		for j := range maze[i] {
			maze[i][j] = true
		}
	}

	type Wall struct {
		col, row   int
		cell1      int
		cell2      int
		isVertical bool
	}

	walls := make([]Wall, 0)

	for row := 1; row < mazeRows-1; row += 2 {
		for col := 1; col < mazeCols-1; col += 2 {
			cellIdx := (row/2)*((mazeCols-1)/2) + (col / 2)

			if col < mazeCols-2 {
				walls = append(walls, Wall{
					col:        col + 1,
					row:        row,
					cell1:      cellIdx,
					cell2:      cellIdx + 1,
					isVertical: true,
				})
			}

			if row < mazeRows-2 {
				walls = append(walls, Wall{
					col:        col,
					row:        row + 1,
					cell1:      cellIdx,
					cell2:      cellIdx + ((mazeCols - 1) / 2),
					isVertical: false,
				})
			}
		}
	}

	for row := 1; row < mazeRows; row += 2 {
		for col := 1; col < mazeCols; col += 2 {
			maze[row][col] = false
		}
	}

	rand.Shuffle(len(walls), func(i, j int) {
		walls[i], walls[j] = walls[j], walls[i]
	})

	ds := newDisjointSet(((mazeRows - 1) / 2) * ((mazeCols - 1) / 2))

	for _, wall := range walls {
		if findFunc(ds, wall.cell1) != findFunc(ds, wall.cell2) {
			maze[wall.row][wall.col] = false
			union(ds, wall.cell1, wall.cell2)
		}
	}
}

func (m *MazeGenerator) DisplayMaze() string {
	if m.maze == nil {
		return ""
	}

	result := ""
	blocks := BlockStyles[m.blockStyle]

	for i := 0; i < len(m.maze); i++ {
		for j := 0; j < len(m.maze[i]); j++ {
			if m.maze[i][j] {
				result += blocks.Wall
			} else {
				result += blocks.Path
			}
		}
		result += "\n"
	}
	return result
}
