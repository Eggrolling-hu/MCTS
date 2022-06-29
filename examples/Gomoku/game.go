package main

import (
	"fmt"
	"mcts"
	"strings"
)

type Move struct {
	X        int `json:"x"`
	Y        int `json:"y"`
	PlayerID int `json:"playerID"`
}

func (m Move) Log() string {
	return fmt.Sprintf("%d -> (%d, %d)", m.PlayerID, m.X, m.Y)
}

type GomokuState struct {
	Board        [][]int `json:"board"`
	BoardSize    int     `json:"boardSize"`
	ChainSize    int     `json:"chainSize"`
	BackMove     Move    `json:"backMove"`
	InitPlayerID int     `json:"initPlayerID"`
	NextPlayerID int     `json:"nextPlayerID"`
}

func (s *GomokuState) GetAvailableActions() []mcts.Action {
	if s.IsForcedTerminated() {
		return []mcts.Action{}
	}
	xRange := []int{s.BoardSize, 0}
	yRange := []int{s.BoardSize, 0}
	// n := int(s.ChainSize / 2.0)
	n := 1

	for x := 0; x < s.BoardSize; x++ {
		for y := 0; y < s.BoardSize; y++ {
			if s.Board[x][y] != 0 {
				xRange[0] = Min(x, xRange[0])
				xRange[1] = Max(x, xRange[1])
				yRange[0] = Min(y, yRange[0])
				yRange[1] = Max(y, yRange[1])
			}
		}
	}

	availableActions := make([]mcts.Action, 0, s.BoardSize*s.BoardSize)

	if xRange[0] == s.BoardSize && xRange[1] == 0 && yRange[0] == s.BoardSize && yRange[1] == 0 {
		center := s.BoardSize/2 + 1
		availableAction := Move{center, center, s.NextPlayerID}
		availableActions = append(availableActions, availableAction)
		fmt.Println(availableActions)
		return availableActions
	}

	for x := 0; x < s.BoardSize; x++ {
		for y := 0; y < s.BoardSize; y++ {
			if x < xRange[0]-n || x > xRange[1]+n {
				continue
			}
			if y < yRange[0]-n || y > yRange[1]+n {
				continue
			}
			if s.Board[x][y] != 0 {
				continue
			}

			x1 := Max(0, x-n)
			x2 := Min(x+n, s.BoardSize-n)
			y1 := Max(0, y-n)
			y2 := Min(y+n, s.BoardSize-n)

			if GetIsAlone(s.Board, x1, x2, y1, y2) {
				continue
			}

			availableAction := Move{x, y, s.NextPlayerID}
			availableActions = append(availableActions, availableAction)
		}
	}

	// tricky: draw's backMove playerID is 0
	// if len(availableActions) == 0 {
	// 	s.BackMove.PlayerID = 0
	// }

	return availableActions
}

func (s *GomokuState) TakeAction(action mcts.Action) {
	p := action.(Move)
	s.BackMove = p
	s.Board[p.X][p.Y] = p.PlayerID
	s.NextPlayerID *= -1
}

func (s *GomokuState) Evaluate() float64 {
	if s.BackMove.PlayerID != s.InitPlayerID {
		return 0.0
	}
	return 1.0
}

func (s *GomokuState) IsForcedTerminated() bool {
	if s.BackMove.PlayerID == 0 {
		return false
	}
	sWin := strings.Repeat(fmt.Sprintf("%+d", s.BackMove.PlayerID), s.ChainSize)
	if strings.Contains(XString(s), sWin) {
		return true
	}
	if strings.Contains(YString(s), sWin) {
		return true
	}
	if strings.Contains(QString(s), sWin) {
		return true
	}
	if strings.Contains(EString(s), sWin) {
		return true
	}

	return false
}

func (s *GomokuState) IsRivalRound() bool {
	return s.InitPlayerID == s.NextPlayerID
}

func (s *GomokuState) Log() string {
	backMoveLog := s.BackMove.Log()
	boardLog := ""
	for x := 0; x < s.BoardSize; x++ {
		for y := 0; y < s.BoardSize; y++ {
			boardLog += fmt.Sprintf("%+d ", s.Board[x][y])
		}
		boardLog += "\n"
	}
	boardLog = boardLog[:len(boardLog)-1]
	return fmt.Sprintf("backMove: %s\nboard:\n%s", backMoveLog, boardLog)
}

func (s *GomokuState) Deepcopy() mcts.State {
	board := make([][]int, 0, len(s.Board))
	for i := range s.Board {
		newLine := make([]int, len(s.Board[i]))
		copy(newLine, s.Board[i])
		board = append(board, newLine)
	}
	backMove := Move{s.BackMove.X, s.BackMove.Y, s.BackMove.PlayerID}
	newGomokuState := GomokuState{
		Board:        board,
		BoardSize:    s.BoardSize,
		ChainSize:    s.ChainSize,
		BackMove:     backMove,
		InitPlayerID: s.InitPlayerID,
		NextPlayerID: s.NextPlayerID,
	}
	return &newGomokuState
}
