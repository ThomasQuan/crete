package main

import (
	"strconv"
)

func main() {
	layout := NewLayout()

	layout.AlgorithmList.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		switch mainText {
		case "DFS":
			if !layout.matrix.isSolving {
				width, height := layout.GetMatrixDimensions()
				layout.LogView.SetText(layout.LogView.GetText(true) + "\nGenerating maze using DFS of size " + strconv.Itoa(width) + "x" + strconv.Itoa(height))
				layout.MatrixBox.SetText(layout.matrix.GenerateMaze(height, width, DFS))
			}
		case "Prim's":
			if !layout.matrix.isSolving {
				width, height := layout.GetMatrixDimensions()
				layout.LogView.SetText(layout.LogView.GetText(true) + "\nGenerating maze using Prim's algorithm of size " + strconv.Itoa(width) + "x" + strconv.Itoa(height))
				layout.MatrixBox.SetText(layout.matrix.GenerateMaze(height, width, Prims))
			}
		case "Kruskal's":
			if !layout.matrix.isSolving {
				width, height := layout.GetMatrixDimensions()
				layout.LogView.SetText(layout.LogView.GetText(true) + "\nGenerating maze using Kruskal's algorithm of size " + strconv.Itoa(width) + "x" + strconv.Itoa(height))
				layout.MatrixBox.SetText(layout.matrix.GenerateMaze(height, width, Kruskals))
			}
		case "Solve Maze (DFS)":
			if !layout.matrix.isSolving {
				layout.matrix.isSolving = true
				layout.LogView.SetText(layout.LogView.GetText(true) + "\nSolving maze using DFS...")
				currentMaze := layout.MatrixBox.GetText(true)
				solver := NewSolver()

				go func() {
					solution := solver.SolveMaze(layout.matrix, currentMaze, DFSSearchType)
					layout.App.QueueUpdateDraw(func() {
						layout.MatrixBox.SetText(solution)
						layout.matrix.isSolving = false
					})
				}()
			}

		case "Solve Maze (BFS)":
			if !layout.matrix.isSolving {
				layout.matrix.isSolving = true
				layout.LogView.SetText(layout.LogView.GetText(true) + "\nSolving maze using BFS...")
				currentMaze := layout.MatrixBox.GetText(true)
				// Prevent deadlock, basically an async await
				go func() {
					solver := NewSolver()
					solution := solver.SolveMaze(layout.matrix, currentMaze, BFSSearchType)
					layout.App.QueueUpdateDraw(func() {
						layout.MatrixBox.SetText(solution)
					})
				}()
			}

		case "Fast Mode":
			if !layout.matrix.isSolving {
				layout.matrix.fastMode = !layout.matrix.fastMode
				if layout.matrix.fastMode {
					layout.AlgorithmList.SetItemText(5, "Fast Mode", "Fast mode is ON")
				} else {
					layout.AlgorithmList.SetItemText(5, "Fast Mode", "Fast mode is OFF")
				}
				layout.LogView.SetText("Fast mode is turned " + strconv.FormatBool(layout.matrix.fastMode))
			}
		}

	})

	if err := layout.Run(); err != nil {
		panic(err)
	}
}
