package game

type Game struct {
	Id               *string `json:"id"`
	Black            *string `json:"black"`
	White            *string `json:"white"`
	Board            Board   `json:"board"`
	InitialRoll      bool    `json:"initialRoll"`
	BlackInitialRoll *int8   `json:"blackInitialRoll"`
	WhiteInitialRoll *int8   `json:"whiteInitialRoll"`
	CurrentRoll      *Roll   `json:"currentRoll"`
	CurrentTurn      *Color  `json:"currentTurn"`
}

type Player struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type Board struct {
	Spaces [24]Triangle `json:"spaces"`
}

type Triangle struct {
	Color Color `json:"color"`
	Count int8  `json:"count"`
}

type Roll [2]int

type Color int

const (
	Black Color = 0
	White Color = 1
	None  Color = -1
)

func NewGame(gameID string) Game {
	return Game{
		Id:               &gameID,
		Black:            nil,
		White:            nil,
		Board:            newBoard(),
		InitialRoll:      true,
		BlackInitialRoll: nil,
		WhiteInitialRoll: nil,
		CurrentRoll:      nil,
		CurrentTurn:      nil,
	}
}

func newBoard() Board {
	triangles := [24]Triangle{
		Triangle{
			Color: Black,
			Count: 2,
		},
		Triangle{
			Color: None,
			Count: 0,
		},
		Triangle{
			Color: None,
			Count: 0,
		},
		Triangle{
			Color: None,
			Count: 0,
		},
		Triangle{
			Color: None,
			Count: 0,
		},
		Triangle{
			Color: White,
			Count: 5,
		},
		Triangle{
			Color: None,
			Count: 0,
		},
		Triangle{
			Color: White,
			Count: 3,
		},
		Triangle{
			Color: None,
			Count: 0,
		},
		Triangle{
			Color: None,
			Count: 0,
		},
		Triangle{
			Color: None,
			Count: 0,
		},
		Triangle{
			Color: Black,
			Count: 5,
		},
		Triangle{
			Color: White,
			Count: 5,
		},
		Triangle{
			Color: None,
			Count: 0,
		},
		Triangle{
			Color: None,
			Count: 0,
		},
		Triangle{
			Color: None,
			Count: 0,
		},
		Triangle{
			Color: Black,
			Count: 3,
		},
		Triangle{
			Color: None,
			Count: 0,
		},
		Triangle{
			Color: Black,
			Count: 5,
		},
		Triangle{
			Color: None,
			Count: 0,
		},
		Triangle{
			Color: None,
			Count: 0,
		},
		Triangle{
			Color: None,
			Count: 0,
		},
		Triangle{
			Color: None,
			Count: 0,
		},
		Triangle{
			Color: White,
			Count: 2,
		},
	}

	return Board{Spaces: triangles}
}
