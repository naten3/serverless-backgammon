package main

import (
	"serverless-backgammon/game"
	// todo better way to import this
)

func StartGame(*game game.Game) *game.Game {
	game.InitialRoll = true
}
