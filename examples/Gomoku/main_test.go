package main

import (
	"fmt"
	"math"
	"mcts"
	"testing"
)

func TestL0(t *testing.T) {

	boardSize := 19
	chainSize := 5
	backMove := Move{}
	initPlayerID := 1
	nextPlayerID := 1

	board := [][]int{}
	for x := 0; x < boardSize; x++ {
		newLine := []int{}
		for y := 0; y < boardSize; y++ {
			newLine = append(newLine, 0)
		}
		board = append(board, newLine)
	}

	s := GomokuState{
		board,
		boardSize,
		chainSize,
		backMove,
		initPlayerID,
		nextPlayerID,
	}

	fmt.Println(s)
	c := s.Deepcopy()
	a := c.GetAvailableActions()
	c.TakeAction(a[0])

	fmt.Println(s)
	fmt.Println(c)

	fmt.Println(s)
	a = s.GetAvailableActions()
	s.TakeAction(a[0])
	fmt.Println(s)
	for i := 0; i < 20; i++ {
		fmt.Println(i)
		a = s.GetAvailableActions()
		s.TakeAction(a[0])
		fmt.Println(s)
		if s.IsForcedTerminated() {
			break
		}
	}
}

func TestL1(t *testing.T) {
	boardSize := 15
	chainSize := 5
	backMove := Move{}
	initPlayerID := 1
	nextPlayerID := 1

	board := [][]int{}
	for x := 0; x < boardSize; x++ {
		newLine := []int{}
		for y := 0; y < boardSize; y++ {
			newLine = append(newLine, 0)
		}
		board = append(board, newLine)
	}

	s := GomokuState{
		board,
		boardSize,
		chainSize,
		backMove,
		initPlayerID,
		nextPlayerID,
	}

	bestAction := mcts.UCTSearch(mcts.State(&s), int64(1e5), int64(5e3), 1/math.Sqrt(2))
	fmt.Println(bestAction)
}

func TestL2(t *testing.T) {
	boardSize := 4
	chainSize := 3
	backMove := Move{2, 2, -1}
	initPlayerID := 1
	nextPlayerID := 1

	board := [][]int{
		{1, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, -1, -1},
		{0, 0, 0, 0},
	}

	s := GomokuState{
		board,
		boardSize,
		chainSize,
		backMove,
		initPlayerID,
		nextPlayerID,
	}

	bestAction := mcts.UCTSearch(mcts.State(&s), int64(5e4), int64(5e2), 1/math.Sqrt(2))
	fmt.Println(bestAction)
	fmt.Println(s.Log())
	s.TakeAction(bestAction)
	fmt.Println(s.Log())
}
