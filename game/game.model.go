package game

type Game struct {
	Id    string
	Black Player
	White Player
	Board Board
}

type Player struct {
	Name string
	Id   string
}

type Board struct {
	spaces [24]Triangle
}

type Triangle struct {
	color Color
	count int8
}

type Color int

const (
	Black Color = 0
	White Color = 1
)
