package main

import (
	"fmt"
	"time"
	"math"
)

type Position struct {
	x, y int
}

var nodesSearched = 0

func getSearchSpace(board *Board, freeThrees *[2]Board, player int) []Position {

	if debug { defer timeFunc(time.Now(), "getSearchSpace") }

	moves := make([]Position, 0, 10)
	alreadyChecked := [19][19]bool {}

	checkAxis := func(x, y, incx, incy int) {
		if isValidMove(board, freeThrees, x + incx, y + incy, player) && !alreadyChecked[y + incy][x + incx] {
			alreadyChecked[y + incy][x + incx] = true
			moves = append(moves, Position{x + incx, y + incy})
		}
		if isValidMove(board, freeThrees, x - incx, y - incy, player) && !alreadyChecked[y - incy][x - incx] {
			alreadyChecked[y - incy][x - incx] = true
			moves = append(moves, Position{x - incx, y - incy})
		}
	}

	for y := 0; y < 19; y++ {
		for x := 0; x < 19; x++ {
			if board[y][x] != empty {
				for i := 0; i < 2; i++ {
					checkAxis(x, y, i, 0)
					checkAxis(x, y, 0, i)
					checkAxis(x, y, i, i)
					checkAxis(x, y, i, -i)
				}
			}
		}
	}

	return moves
}

func search(values *Board, freeThree *[2]Board, player, x, y, depth int, capture *[3]int) (int, int, BoardData) {

	nodesSearched = 0
	startTime := time.Now()

	var	ax, ay int
	var boardData BoardData

	moves := getSearchSpace(values, freeThree, player)

	bestscore := math.MinInt32
	alpha := math.MinInt32
	beta := math.MaxInt32

	// TODO: Parallelize

	for _, pos := range(moves) {
		b := AIBoard{*values, *freeThree, *capture, player}
		move := b.CreateMove(pos)
		b.DoMove(move)
		boardData[move.pos.y][move.pos.x][6] = 1
		score := evaluateBoard(b.Board(), move.pos.x, move.pos.y, player, &boardData, b.CapturesNb())
		boardData[move.pos.y][move.pos.x][5] = score
		boardData[move.pos.y][move.pos.x][7] = score
		if score >= 2e9 {
			b.UndoMove(move)
			return move.pos.x, move.pos.y, boardData
		}
		b.UpdateFreeThrees(move.pos, move.captures)
		s := -searchdeeper(&b, move.pos, depth - 1, -beta, -alpha)
		boardData[move.pos.y][move.pos.x][5] = s
		b.UndoMove(move)
		b.UpdateFreeThrees(move.pos, move.captures)

		if s >= beta {
			return move.pos.x, move.pos.y, boardData
		}

		if s > bestscore {
			bestscore = s
			ax, ay = move.pos.x, move.pos.y
			if s > alpha {
				alpha = s
			}
		}
	}

	fmt.Println(nodesSearched, "nodes searched in", time.Since(startTime), "(", time.Since(startTime) / time.Duration(nodesSearched), "by node)")
	fmt.Println("BEST : ", bestscore)

	return ax, ay, boardData
}

func searchdeeper(b *AIBoard, move Position, depth int, alpha, beta int) int {

	nodesSearched++
	bestscore := math.MinInt32

	if depth == 0 {
		return b.Evaluate(move)
	}

	// Why not ?
	b.SwitchPlayer()
	defer b.SwitchPlayer()

	moves := b.GetNextMoves()

	for _, move := range(moves) {
		if move.Score() >= 2e9 {
			return move.Score()
		}
		b.DoMove(move)
		s := -searchdeeper(b, move.Position(), depth - 1, -beta, -alpha)
		b.UndoMove(move)

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

	if debug { defer timeFunc(time.Now(), "doMove") }

	board[y][x] = player
	getCaptures(board, x, y, player, captures)
	doCaptures(board, captures)
}

func undoMove(board *Board, x, y, player int, captures *[]Position) {

	if debug { defer timeFunc(time.Now(), "undoMove") }

	board[y][x] = empty
	undoCaptures(board, captures, player)
}

func checkAlign(values *Board, x, y, player int) int {
	f := func (incx, incy, x, y, p int) int {
		cnt := 0
		x, y = x + incx, y + incy
		for i := 0; i < 4; i++ {
			if !isInBounds(x, y) || values[y][x] == -p{
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
		if !isInBounds(x + 3 * incx, y + 3 * incy) {
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

// TODO: Deep search doesnt need a BoardData
func evaluateBoard(values *Board, x, y, player int, copy *BoardData, capture *[3]int) int {
	// C'est de la grosse merde !
	// -v2

	if debug { defer timeFunc(time.Now(), "evaluateBoard") }

	var v1, v2 int

	v1 = checkAlign(values, x, y, player)
	copy[y][x][0] = v1
	if v1 >= 4 || capture[player + 1] >= 10 {
		return math.MaxInt32
	}
	v2 = checkAlign(values, x, y, -player)
	copy[y][x][1] = v2

	return v1 + v2 * 2 + capture[player + 1] - capture[-player + 1]
}
