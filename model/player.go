package model

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"math/rand"
	"time"
)

type Player struct {
	Id         int
	Name       string
	Colour     string
	CurrentPos int
}

type PlayerActivity interface {
	Move(Board) error
	Roll(totalPlayers int, board Board, startingVal int) int // Return a positional increment and also consider that player will roll again if he gets a six.
	GetCurrentPos() int
	SetCurrentPos(int)
}

func (player *Player) Move(board Board) error {

	gameStatus := (board).GetGameStatus()
	if gameStatus != INPROGRESS {
		return nil
	}
	boardType := board.GetGameType()

	if boardType == TicTacToe {
		var currPos int
		fmt.Print("Play your Move")

		for {
			fmt.Scanf("%d", &currPos)
			if (currPos < 9) && (currPos > -1) {
				break
			}
			fmt.Print("Invalid Move. Play Again.")
		}

		// Positional validation needs to be done. One position cannot be played again
		board.(*TicTacToeBoard).SetPositionValue(currPos, board.GetPlayers()[board.GetPossession()].Id)
		if board.(*TicTacToeBoard).CheckWinner() {
			log.Infof("Winner Found for TicTacToe game %v ", board.GetPlayers()[board.GetPossession()])
			board.GetChannel() <- &board.GetPlayers()[board.GetPossession()]
		}
		board.(*TicTacToeBoard).ChangePossession()

	} else if boardType == SnakesAndLadder {
		rolledVal := player.Roll((board).GetPlayerCount(), board, 0)

		if player.CurrentPos+rolledVal == 100 {
			board.GetChannel() <- player
		}

		newPos := player.CurrentPos + rolledVal
		if (board).IsSpecialPosition(newPos) {
			typePos := (board).GetSpecialPosition(newPos).Type
			newPos = (board).GetSpecialPosition(newPos).End
			log.Infof("Caught at special position %v Player %v , OldPos %v Newpos %v", typePos, player.Name, player.CurrentPos, newPos)
		}
		player.SetCurrentPos(newPos)
		log.Infof("Player %v moved to position %v", player.Name, player.CurrentPos)

	}

	return nil
}

func ChangePossesion(board Board, totalCount int) {
	newPossession := ((board).GetPossession() + 1) % totalCount
	(board).SetPossession(newPossession)
}

func RollTheDie() int {
	rand.Seed(time.Now().UnixNano())
	rollVal := rand.Intn(6) + 1

	return rollVal
}

func (player *Player) SetCurrentPos(pos int) {
	player.CurrentPos = pos
}
func (player *Player) Roll(totalCount int, board Board, startingVal int) int {

	// Roll till we get a value less than 6
	// If at any point the currentval + rollval  + currentPOs > 100 we return and change possession

	roll := RollTheDie()
	limitFlag := false
	for {
		if roll == 6 {
			startingVal += roll
			if startingVal+player.CurrentPos > 100 {
				limitFlag = true
				break
			}
			roll = RollTheDie()
		} else {
			startingVal += roll
			if startingVal+player.CurrentPos > 100 {
				limitFlag = true
			}
			break
		}
	}

	ChangePossesion(board, totalCount)

	if limitFlag {
		return 0
	}
	return startingVal
}

func (player *Player) GetCurrentPos() int {
	if player != nil {
		return 0
	}
	return player.CurrentPos
}

func GetAllPlayers() []Player {
	playerList := []Player{
		{1, "Vinay", "Blue", 0},
		{2, "Swati", "Green", 0},
	}
	return playerList
}

func ChooseFirst() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2)
}
