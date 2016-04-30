package main

import (
	"fmt"
	"time"
)

type Position struct {
	x, y int
}

var nodesSearched = 0

// copy[0] "score" ia
// copy[1] "score" player 
// copy[2] "capturable"
// copy[3] forbiden ia
// copy[4] forbiden player

// var stopByTime = false
// var node = 0

func bitCount(b int) int {
		var j int
		for i := uint(0); i < 4; i++ {
			if b & (1 << i) != 0 { j++}
		}
		return j
}

func getSearchSpace(board *Board, freeThree *[2]Board, player int) map[Position]bool {

	defer timeFunc(time.Now(), "getSearchSpace")

	m := make(map[Position]bool)

	playerId := (player + 1) / 2
	for i:= 0;i < 19; i++ {
		for j:= 0;j < 19; j++ {
			if board[i][j] != 0 {

				for circle := 1; circle < 3; circle++ {
					a, b := i - circle, i + circle
					for incx := j - (circle - 1); incx < j + circle; incx++ {
						if incx < 0 { incx = 0 } else if incx > 18 { break }
						if a >= 0 && board[a][incx] == 0 && bitCount(freeThree[playerId][a][incx]) != 2 {
							m[Position{incx, a}] = true
						}
						if b < 19 && board[b][incx] == 0 && bitCount(freeThree[playerId][b][incx]) != 2 {
							m[Position{incx, b}] = true
						}
					}
					a, b = j - circle, j + circle
					for incy := i - circle; incy <= i + circle; incy++ {
						if incy < 0 { incy = 0 } else if incy > 18 { break }
						if a >= 0 && board[incy][a] == 0 && bitCount(freeThree[playerId][incy][a]) != 2 {
							m[Position{a, incy}] = true
						}
						if b < 19 && board[incy][b] == 0 && bitCount(freeThree[playerId][incy][b]) != 2 {
							m[Position{b, incy}] = true
						}
					}
				}
			}
		}
	}
	return m
}

func search(values *Board, freeThree *[2]Board, player, x, y, depth int, capture *[3]int) (int, int, BoardData) {

	nodesSearched = 0

	var	ax, ay int
	var copy BoardData

	m := getSearchSpace(values, freeThree, player)

	bestscore := -int(^uint32(0)>>1) // int le plus large possible dans les negatifs

	alpha := -int(^uint32(0)>>1)
	beta := int(^uint32(0)>>1)

	for i, _ := range(m) {
		score := evaluateBoard(values, i.x, i.y, player, &copy, capture)
		if score >= 20 {
			return i.x, i.y, copy
		}

		captures := make([]Position, 0, 16)
		doMove(values, i.x, i.y, player, &captures)
		capture[player + 1] += len(captures)
		s := -searchdeeper(values, freeThree, -player, i.x, i.y, depth - 1, capture, -beta, -alpha)
		undoMove(values, i.x, i.y, player, &captures)
		capture[player + 1] -= len(captures)

		if s >= beta {
			return i.x, i.y, copy
		}
		if s > bestscore {
			bestscore = s
			ax, ay = i.x, i.y
			if s > alpha {
				alpha = s
			}
		}
	}
	fmt.Println(nodesSearched)

	return ax, ay, copy
}

func searchdeeper(values *Board, freeThree *[2]Board, player, x, y, depth int, capture *[3]int, alpha, beta int) int {

	//var copy_capt Board

	nodesSearched++

	var copy BoardData

	if depth == 0 {
		return evaluateBoard(values, x, y, player, &copy, capture)
	}

	m := getSearchSpace(values, freeThree, player)

	bestscore := -int(^uint32(0)>>1) // int le plus large possible dans les negatif

	for i, _ := range(m) {
		score := evaluateBoard(values, i.x, i.y, player, &copy, capture)
		if score >= 20 {
			return score
		}

		captures := make([]Position, 0, 16)
		doMove(values, i.x, i.y, player, &captures)
		capture[player + 1] += len(captures)
		updateFreeThrees(values, freeThree, x, y, player, captures)
		s := -searchdeeper(values, freeThree, -player, i.x, i.y, depth - 1, capture, -beta, -alpha)
		undoMove(values, i.x, i.y, player, &captures)
		capture[player+1] -= len(captures)
		updateFreeThrees(values, freeThree, x, y, player, captures)

		if s >= beta {
			return s
		}
		if s > bestscore {
			bestscore = s
			if s > alpha {
				alpha = s
			}
		}
	}
	return bestscore
}

func doMove(board *Board, x, y, player int, captures *[]Position) {
	defer timeFunc(time.Now(), "doMove")

	board[y][x] = player
	getCaptures(board, x, y, player, captures)
	doCaptures2(board, captures)
}

func undoMove(board *Board, x, y, player int, captures *[]Position) {
	defer timeFunc(time.Now(), "undoMove")

	board[y][x] = empty
	undoCaptures(board, captures, player)
}

func checkAlign(values *Board, x, y, player int) int {
	f := func (incx, incy, x, y, p int) int {
		cnt := 0
		x, y = x + incx, y + incy
		for i := 0; i < 4; i++ {
			if !checkBounds(x, y) || values[y][x] == -p{
				return cnt
			}
			if values[y][x] == p {
				cnt += 1
			}
			x += incx
			y += incy
		}
		return cnt
	}
	max, t := 0, 0 
	
	max = f(-1, -1, x, y, player) + f(1, 1, x, y, player)
	t = f(1, -1, x, y, player) + f(-1, 1, x, y, player)
	if t > max {
		max = t
	}	
	t = f(0, -1, x, y, player) + f(0, 1, x, y, player)
	if t > max {
		max = t
	}
	t = f(-1, 0, x, y, player) + f(1, 0, x, y, player)
	if t > max {
		return t
	} 
	return max
}

func checkCapt(values *Board, x, y, player int) int {
	capt := func (incx, incy int) int {
		if !checkBounds(x + 3 * incx, y + 3 * incy) {
			return 0
		}
		if	values[y + incy][x + incx] == -player &&
			values[y + 2 * incy][x + 2 * incx] == -player &&
		 	values[y + 3 * incy][x + 3 * incx] == player {
				return 2
		}
		return 0
	}
	return  capt(-1, -1) + capt(1, 1) + capt(1, -1) + capt(-1, 1) +
			capt(0, -1) + capt(0, 1) + capt(-1, 0) + capt(1, 0)
}

func evaluateBoard(values *Board, x, y, player int, copy *BoardData, capture *[3]int) int {
	
	var v1, v2, v3 int

	v1 = checkAlign(values, x, y, player)
	copy[y][x][0] = v1
	if v1 >= 4 {
		return 20
	}
	v2 = -checkAlign(values, x, y, -player )
	copy[y][x][1] = -v2
	if v2 <= -4 {
		return -20
	}
	v3 = checkCapt(values, x, y, player)
	copy[y][x][2] = v3
	if v3 > 0 {
		if capture[player + 1] + v3 >= 10 {
			return 20
		}
		return capture[player + 1] + v3 + 2
	}
	if (-v2 * 2) > v1 {
		return v2
	}
	return v1
}
