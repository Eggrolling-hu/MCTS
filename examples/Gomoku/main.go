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

func main() {
	Server()
}
