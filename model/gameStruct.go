package model

import (
	"errors"
)

var endGame chan int

type Board interface {
	//ConstructGameBoard() error
	MarkSpecialPositions([]Position) error
	GetSpecialPosition(int) Position
	IsSpecialPosition(position int) bool
	GetPlayers() []Player
	GetPossession() int
	SetPossession(possession int)
	GetPlayerCount() int
	GetGameStatus() GameStatus
	SetGameStatus(status GameStatus)

	Play()
	GetChannel() chan *Player
	//EndGame()
}

type GameBoard struct {
	SpecialPosition map[int]Position
	Players         []Player
	Possession      int
	Status          GameStatus
	Channel         chan *Player
}

type (
	CellType   string
	GameStatus string
)

const (
	SNAKE    = CellType("SNAKE")
	LADDER   = CellType("LADDER")
	ORDINARY = CellType("ORDINARY")

	INPROGRESS = GameStatus("PROGRESS")
	ENDED      = GameStatus("ENDED")
	NOTSTARTED = GameStatus("NOTSTARTED")
)

type Position struct {
	Initial int
	End     int
	Type    CellType
}

func (gameBoard *GameBoard) IsSpecialPosition(pos int) bool {
	if _, ok := gameBoard.SpecialPosition[pos]; !ok {
		return false
	} else {
		return true
	}
}

func (gameBoard *GameBoard) GetSpecialPosition(pos int) Position {
	if val, ok := gameBoard.SpecialPosition[pos]; ok {
		return val
	}
	return Position{}
}

func (gameBoard *GameBoard) GetPlayers() []Player {
	return gameBoard.Players
}

// Assumption is that there is no recursive cycling here. For instance no snake bite can lead to a ladder start which takes it to the same snake position.
func (gameBoard *GameBoard) MarkSpecialPositions(position []Position) error {
	for _, pos := range position {
		if pos.Initial == 100 {
			gameBoard.SpecialPosition = make(map[int]Position)
			return errors.New("Cannot have a special position at final cell.")
		}
		if _, ok := gameBoard.SpecialPosition[pos.Initial]; ok {
			gameBoard.SpecialPosition = make(map[int]Position)
			return errors.New("Cannot have multiple special positions on same cell")
		}
		gameBoard.SpecialPosition[pos.Initial] = pos
	}
	return nil
}

func (gameBoard *GameBoard) GetPossession() int {
	return gameBoard.Possession
}

func (gameBoard *GameBoard) SetPossession(possession int) {
	gameBoard.Possession = possession
}

func (gameBoard *GameBoard) GetPlayerCount() int {
	return len(gameBoard.Players)
}

func (gameBoard *GameBoard) Play() {
	gameBoard.Status = INPROGRESS

	for {
		select {
		case <-endGame:
			break
		default:
			possessionPlayerIndex := gameBoard.GetPossession()
			player := &(*gameBoard).Players[possessionPlayerIndex]
			player.Move(gameBoard)
		}
	}
}

func (gameBoard *GameBoard) SetGameStatus(gameStatus GameStatus) {
	gameBoard.Status = gameStatus
}

func (gameBoard *GameBoard) GetGameStatus() GameStatus {
	return (*gameBoard).Status
}

func (gameBoard *GameBoard) GetChannel() chan *Player {
	return (*gameBoard).Channel
}
func PrepareGameBoard(ladderPos, snakePos []int, players []Player, channel chan *Player) Board {
	gameBoard := GameBoard{
		Players:    players,
		Possession: 0,
		Status:     NOTSTARTED,
		Channel:    channel,
		SpecialPosition: make(map[int]Position),
	}
	specialPositions := []Position{}

	for idx := 0; idx < len(ladderPos); idx += 2 {
		start := ladderPos[idx]
		end := ladderPos[idx+1]
		if end <= start {
			continue
		}
		specialPositions = append(specialPositions, Position{
			Initial: start,
			End:     end,
			Type:    LADDER,
		})
	}
	for idx := 0; idx < len(snakePos); idx += 2 {
		start := snakePos[idx]
		end := snakePos[idx+1]
		if end >= start {
			continue
		}
		specialPositions = append(specialPositions, Position{
			Initial: start,
			End:     end,
			Type:    SNAKE,
		})
	}
	if err := gameBoard.MarkSpecialPositions(specialPositions); err != nil {
		return nil
	}

	return &gameBoard
}
