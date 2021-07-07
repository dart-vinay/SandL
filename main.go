package main

import (
	"SandL/model"
	"github.com/labstack/gommon/log"
)

// How the game actually looks like
// 1. There is a board that count from 1-100
// 2. We have snakes that can take you from pos a to pos b where a>b
// 3. We have ladders that take you from pos c to pos d where d>c
// - On next layer we can specify that ladders can only be vertical, i.e. we will divide the board into 10 equal columns. or other simple way could be to say that
//   we cannot have a ladder with a different between its two pos less than 10.
// 4.

// Methods that we can write
// 1. Create a snake & ladder board.
// - Specify the positions of snakes and their corresponding outcome if person lands on that position. Do the same for all the ladders.
// - Decide on the structure of this positional mapping : Can make map of all the special position.

var winnerChannel chan *model.Player
var endGame chan bool

func InitializeRelevantChannels() {
	winnerChannel = make(chan *model.Player, 1)
	endGame = make(chan bool, 1)
}
func main() {

	log.Infof("Welcome to a game of Snakes and Ladders with Swati & Vinay ")
	var winner *model.Player

	ladderPos := []int{5, 50, 15, 97}
	snakePos := []int{51, 0, 99, 11, 34, 6, 89, 10, 71, 21}

	players := model.GetAllPlayers()

	InitializeRelevantChannels()

	board := model.PrepareGameBoard(ladderPos, snakePos, players, winnerChannel)
	go board.Play()

	winner = <-winnerChannel
	endGame <- true
	board.SetGameStatus(model.ENDED)
	log.Infof("Got a winner %v", winner.Name)
	return

}
