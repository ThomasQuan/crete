package main

type BlockStyle int

const (
	FullBlock BlockStyle = iota
	HalfBlock
	QuarterBlock
	ShadeBlock
)

type Block struct {
	Wall     string
	Path     string
	Solution string
}

var BlockStyles = map[BlockStyle]Block{
	FullBlock: {
		Wall:     "██",
		Path:     "  ",
		Solution: "▒▒",
	},
	HalfBlock: {
		Wall:     "▀▀",
		Path:     "  ",
		Solution: "▄▄",
	},
	QuarterBlock: {
		Wall:     "▌▐",
		Path:     "  ",
		Solution: "░░",
	},
	ShadeBlock: {
		Wall:     "▓▓",
		Path:     "  ",
		Solution: "░░",
	},
}
