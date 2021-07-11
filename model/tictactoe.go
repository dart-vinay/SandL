package model

import (
	"fmt"
	"time"
)

// Need Following Validations:
//1. If game ends in a draw, then how to stop it. Keep a track of the move counts
//2. Valid Moves validation

type TicTacToeBoard struct {
	AllPositions   []int
	Players        []Player
	Possession     int
	Status         GameStatus
	WinnerChannel  chan *Player
	EndGameChannel chan bool

	Type string
}

func (ticTacToeBoard *TicTacToeBoard) MarkSpecialPositions(pos []Position) error {
	return nil
}

func (ticTacToeBoard *TicTacToeBoard) GetSpecialPosition(index int) Position {
	return Position{
		Value: ticTacToeBoard.AllPositions[index],
	}
}

func (ticTacToe *TicTacToeBoard) IsSpecialPosition(pos int) bool {
	return true
}

func (ticTacToeBoard *TicTacToeBoard) GetPlayers() []Player {
	return ticTacToeBoard.Players
}

func (ticTacToeBoard *TicTacToeBoard) GetPossession() int {
	return ticTacToeBoard.Possession
}

func (ticTacToeBoard *TicTacToeBoard) SetPossession(poss int) {
	ticTacToeBoard.Possession = poss
}

func (ticTacToeBoard *TicTacToeBoard) GetPlayerCount() int {

	return len(ticTacToeBoard.Players)
}

func (ticTacToeBoard *TicTacToeBoard) GetGameStatus() GameStatus {
	return ticTacToeBoard.Status
}

func (ticTacToeBoard *TicTacToeBoard) GetGameType() string {
	return ticTacToeBoard.Type
}

func (ticTacToeBoard *TicTacToeBoard) SetGameStatus(status GameStatus) {
	ticTacToeBoard.Status = status
}

func (ticTacToeBoard *TicTacToeBoard) Play() {

	ticTacToeBoard.Status = INPROGRESS
	for {
		select {
		case <-ticTacToeBoard.EndGameChannel:
			ticTacToeBoard.Status = ENDED
			break
		default:
			time.Sleep(10 * time.Millisecond)
			ticTacToeBoard.PrintCurrentGameState()
			currPlayer := &ticTacToeBoard.Players[ticTacToeBoard.GetPossession()]
			currPlayer.Move(ticTacToeBoard)
		}
	}
}

func (ticTacToeBoard *TicTacToeBoard) GetChannel() chan *Player {
	return ticTacToeBoard.WinnerChannel
}

func (ticTacToeBoard *TicTacToeBoard) GetEndGameChannel() chan bool {
	return ticTacToeBoard.EndGameChannel
}

func PrepareTicTacToeBoard(players []Player, winnerChan chan *Player, endgameChan chan bool) Board {
	board := &TicTacToeBoard{
		AllPositions:   []int{0, 0, 0, 0, 0, 0, 0, 0, 0},
		Players:        players,
		WinnerChannel:  winnerChan,
		EndGameChannel: endgameChan,
		Possession:     0,
		Status:         NOTSTARTED,
		Type:           "TicTacToe",
	}

	return board
}

func (ticTacToeBoard *TicTacToeBoard) PrintCurrentGameState() {
	idx := 0

	for ; idx < 9; idx += 3 {
		fmt.Println(ticTacToeBoard.AllPositions[idx], ticTacToeBoard.AllPositions[idx+1], ticTacToeBoard.AllPositions[idx+2])
	}

}

func (ticTacToeBoard *TicTacToeBoard) ChangePossession() {
	ticTacToeBoard.SetPossession(1 - ticTacToeBoard.GetPossession())
}

func (ticTacToeBoard *TicTacToeBoard) CheckWinner() bool {
	potentialWinner := ticTacToeBoard.Players[ticTacToeBoard.GetPossession()].Id
	idx := 0
	//CheckVertical
	for ; idx < 3; idx += 1 {
		if ticTacToeBoard.AllPositions[idx] == potentialWinner && ticTacToeBoard.AllPositions[idx+3] == potentialWinner && ticTacToeBoard.AllPositions[idx+6] == potentialWinner {
			return true
		}
	}

	//CheckHorizontal
	idx = 0
	for ; idx < 9; idx += 3 {
		if ticTacToeBoard.AllPositions[idx] == potentialWinner && ticTacToeBoard.AllPositions[idx+1] == potentialWinner && ticTacToeBoard.AllPositions[idx+2] == potentialWinner {
			return true
		}
	}

	// CheckCrosses
	if (ticTacToeBoard.AllPositions[0] == potentialWinner && ticTacToeBoard.AllPositions[4] == potentialWinner && ticTacToeBoard.AllPositions[8] == potentialWinner) || (ticTacToeBoard.AllPositions[2] == potentialWinner && ticTacToeBoard.AllPositions[4] == potentialWinner && ticTacToeBoard.AllPositions[6] == potentialWinner) {
		return true
	}

	return false
}

func (ticTacToeBoard *TicTacToeBoard) SetPositionValue(pos, val int) {
	ticTacToeBoard.AllPositions[pos] = val
}
