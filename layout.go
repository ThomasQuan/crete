package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Layout struct {
	App           *tview.Application
	MatrixBox     *tview.TextView
	AlgorithmList *tview.List
	LogView       *tview.TextView
	matrix        *MazeGenerator
}

// Basically a state manager for the application
func NewLayout() *Layout {
	app := tview.NewApplication()
	matrixBox := tview.NewTextView()
	logView := tview.NewTextView()

	l := &Layout{
		App:           app,
		MatrixBox:     matrixBox,
		AlgorithmList: tview.NewList(),
		LogView:       logView,

		matrix: &MazeGenerator{
			blockStyle: FullBlock,
			app:        app,
			matrixBox:  matrixBox,
			logView:    logView,
			logger:     NewLogManager(app, logView),
		},
	}
	l.setupLayout()
	return l
}

// Basically a constructor for the layout
func (l *Layout) setupLayout() {
	// Algo list
	l.AlgorithmList.
		AddItem("DFS", "Generate a maze using DFS", '1', nil).
		AddItem("Prim's", "Generate a maze using Prim's algorithm", '2', nil).
		AddItem("Kruskal's", "Generate a maze using Kruskal's algorithm", '3', nil).
		AddItem("Solve Maze (DFS)", "Solve the maze using DFS", 'd', nil).
		AddItem("Solve Maze (BFS)", "Solve the maze using BFS", 'b', nil).
		AddItem("Fast Mode", "Fast mode is OFF", 'f', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			l.App.Stop()
		})

	l.AlgorithmList.SetBorder(true).SetTitle("Select Algorithm")
	l.AlgorithmList.SetMouseCapture(nil)

	// Log list
	l.LogView.
		SetChangedFunc(func() {
			l.App.Draw()
			l.LogView.ScrollToEnd()
		}).
		SetDynamicColors(true).
		SetBorder(true).
		SetTitle("Log")

	// Matrix display
	l.MatrixBox.
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
			heightReduction := int(float64(height) * 0.1)
			y += heightReduction
			height -= heightReduction
			return x, y, width, height
		}).
		SetBorder(true).
		SetTitle("Matrix Display")

	// Arrangement
	rightSide := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(l.AlgorithmList, 0, 1, true).
		AddItem(l.LogView, 0, 1, false)

	flex := tview.NewFlex().
		AddItem(l.MatrixBox, 0, 7, false).
		AddItem(rightSide, 0, 3, true)

	l.App.SetRoot(flex, true).EnableMouse(true)
}

func (l *Layout) Run() error {
	return l.App.Run()
}

// Get the dimensions of the matrix box
func (l *Layout) GetMatrixDimensions() (width, height int) {
	_, _, width, height = l.MatrixBox.GetRect()
	heightReduction := int(float64(height) * 0.1)
	height -= heightReduction * 2
	width = (width - 10) / 2
	return
}
