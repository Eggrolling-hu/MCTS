package main

import "fmt"

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func XString(s *GomokuState) string {
	chain := ""
	xMin := Max(0, s.BackMove.X-s.ChainSize+1)
	xMax := Min(s.BackMove.X+s.ChainSize, s.BoardSize)
	for x := xMin; x < xMax; x++ {
		chain += fmt.Sprintf("%+d", s.Board[x][s.BackMove.Y])
	}
	return chain
}

func YString(s *GomokuState) string {
	chain := ""
	yMin := Max(0, s.BackMove.Y-s.ChainSize+1)
	yMax := Min(s.BackMove.Y+s.ChainSize, s.BoardSize)
	for y := yMin; y < yMax; y++ {
		chain += fmt.Sprintf("%+d", s.Board[s.BackMove.X][y])
	}
	return chain
}

func QString(s *GomokuState) string {
	chain := ""
	for e := -1 * (s.ChainSize - 1); e < s.ChainSize; e++ {
		x := s.BackMove.X + e
		y := s.BackMove.Y + e
		if x < 0 || x >= s.BoardSize {
			continue
		}
		if y < 0 || y >= s.BoardSize {
			continue
		}
		chain += fmt.Sprintf("%+d", s.Board[x][y])
	}
	return chain
}

func EString(s *GomokuState) string {
	chain := ""
	for e := -1 * (s.ChainSize - 1); e < s.ChainSize; e++ {
		x := s.BackMove.X + e
		y := s.BackMove.Y - e
		if x < 0 || x >= s.BoardSize {
			continue
		}
		if y < 0 || y >= s.BoardSize {
			continue
		}
		chain += fmt.Sprintf("%+d", s.Board[x][y])
	}
	return chain
}

func GetIsAlone(board [][]int, x1, x2, y1, y2 int) bool {
	for i := x1; i <= x2; i++ {
		for j := y1; j <= y2; j++ {
			if board[i][j] != 0 {
				return false
			}
		}
	}
	return true
}
