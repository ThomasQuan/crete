package main

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
)

type LogManager struct {
	view *tview.TextView
	app  *tview.Application
}

func NewLogManager(app *tview.Application, view *tview.TextView) *LogManager {
	return &LogManager{
		view: view,
		app:  app,
	}
}

// Log a message to the log view
func (l *LogManager) Log(msg string) {
	timestamp := time.Now().Format("15:04:05")
	formattedMsg := fmt.Sprintf("[yellow]%s[white] %s", timestamp, msg)
	if l.app != nil {
		l.app.QueueUpdateDraw(func() {
			l.view.SetText(l.view.GetText(true) + formattedMsg)
		})
	}
}

// Formatted Log Step 
func (l *LogManager) LogStep(stepCount int, x, y int, isBacktracking bool) {
	direction := "exploring"
	if isBacktracking {
		direction = "backtracking"
	}
	msg := fmt.Sprintf("[blue]Step %d: %s at (%d,%d)[white]\n",
		stepCount, direction, x, y)
	l.Log(msg)
}

// Formatted Log Start
func (l *LogManager) LogStart(algorithm string) {
	l.Log(fmt.Sprintf("[green]Starting %s maze solution...[white]\n", algorithm))
}

// Formatted Log Complete
func (l *LogManager) LogComplete(stepCount, backtrackCount int) {
	l.Log(fmt.Sprintf("[green]Finished! Total steps: %d, Backtracks: %d[white]\n",
		stepCount, backtrackCount))
}

// Formatted Log No Solution
func (l *LogManager) LogNoSolution() {
	l.Log("[red]No solution found![white]\n")
}

// Formatted Log Reset
func (l *LogManager) Reset() {
	l.view.SetText("")
}
