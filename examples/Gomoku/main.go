package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"mcts"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func Test() {

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

func StandAlone(s mcts.State, ch chan Move, wg *sync.WaitGroup) {
	bestAction := mcts.UCTSearch(s, int64(6e4), int64(3e2), 1/math.Sqrt(2))
	move := bestAction.(Move)

	ch <- move
	wg.Done()
}

func Server() {
	router := gin.Default()
	router.POST("/Gomoku", func(c *gin.Context) {
		var s GomokuState
		if err := c.ShouldBindJSON(&s); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("\033[1;32mreceive new data from %s \033[0m", c.ClientIP())

		maxNumProcess := 5
		var wg sync.WaitGroup
		wg.Add(maxNumProcess)
		ch := make(chan Move, maxNumProcess)

		for i := 0; i < maxNumProcess; i++ {
			go StandAlone(s.Deepcopy(), ch, &wg)
		}
		wg.Wait()
		close(ch)

		moves := []Move{}
		for m := range ch {
			moves = append(moves, m)
		}
		fmt.Println(moves)
		rand.Seed(time.Now().UnixNano())
		randIdx := rand.Intn(len(moves))
		move := moves[randIdx]

		c.JSON(http.StatusOK, gin.H{"move": move})
	})

	router.Run(":6666")
}

func Local() {
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

func Local2() {
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

func main() {
	// Local()
	Server()
	// Local2()
}
