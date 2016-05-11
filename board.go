package main

import (
	"time"
	"math"
)

type AIBoard struct {
	board		Board
	// TODO: Split arrays into named members
	freeThrees	[2]Board
	capturesNb	[3]int
	player		int
}

type Move struct {
	pos			Position
	captures	[]Position
}

func (board *AIBoard) Board() *Board {
	return &board.board
}

func (board *AIBoard) FreeThrees() *[2]Board {
	return &board.freeThrees
}

func (board *AIBoard) CapturesNb() *[3]int {
	return &board.capturesNb
}

func (board *AIBoard) Player() *int {
	return &board.player
}

type FreeThreesUpdate struct {
	freeThrees	bool // Dunno. New free threes positions/axes.
}

func (board *AIBoard) SwitchPlayer() {
	board.player = -board.player
}

func (board *AIBoard) isValidMove(x, y int) bool {
	return isInBounds(x, y) &&
		board.board[y][x] == empty &&
		!doesDoubleFreeThree(&board.freeThrees, x, y, board.player)
}

func (board *AIBoard) GetSearchSpace() []Position {

	defer timeFunc(time.Now(), "getSearchSpace")

	moves := make([]Position, 0, 10)
	alreadyChecked := [19][19]bool {}

	checkAxis := func(x, y, incx, incy int) {
		if board.isValidMove(x + incx, y + incy) && !alreadyChecked[y + incy][x + incx] {
			alreadyChecked[y + incy][x + incx] = true
			moves = append(moves, Position{x + incx, y + incy})
		}
		if board.isValidMove(x - incx, y - incy) && !alreadyChecked[y - incy][x - incx] {
			alreadyChecked[y - incy][x - incx] = true
			moves = append(moves, Position{x - incx, y - incy})
		}
	}

	for y := 0; y < 19; y++ {
		for x := 0; x < 19; x++ {
			if board.board[y][x] != empty {
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

func (board *AIBoard) DoMove(pos Position) []Position {
	captures := make([]Position, 0, 16)
	doMove(&board.board, pos.x, pos.y, board.player, &captures)
	board.capturesNb[board.player + 1] += len(captures)
	return captures
}

func (board *AIBoard) UndoMove(pos Position, captures *[]Position) {
	undoMove(&board.board, pos.x, pos.y, board.player, captures)
	board.capturesNb[board.player + 1] -= len(*captures)
}

func (board *AIBoard) UpdateFreeThrees(pos Position, captures []Position) {
	// TODO: Two functions -> update from move and update from move cancelation
	defer timeFunc(time.Now(), "updateFreeThree")

	if (board.board[pos.y][pos.x] != empty) {
		board.freeThrees[0][pos.y][pos.x] = 0
		board.freeThrees[1][pos.y][pos.x] = 0
	}
	checkDoubleThree(&board.board, &board.freeThrees[(board.player + 1) / 2], pos.x, pos.y, board.player)
	checkDoubleThree(&board.board, &board.freeThrees[(-board.player + 1) / 2], pos.x, pos.y, -board.player)
	for _, pos := range captures {
		if (board.board[pos.y][pos.x] != empty) {
			board.freeThrees[0][pos.y][pos.x] = 0
			board.freeThrees[1][pos.y][pos.x] = 0
		}
		checkDoubleThree(&board.board, &board.freeThrees[(board.player + 1) / 2], pos.x, pos.y, board.player)
		checkDoubleThree(&board.board, &board.freeThrees[(-board.player + 1) / 2], pos.x, pos.y, -board.player)
	}
}

func (board *AIBoard) Evaluate(pos Position) int {
	// C'est de la grosse merde !
	// -v2

	defer timeFunc(time.Now(), "evaluateBoard")

	var v1, v2 int

	v1 = board.checkAlign(pos, board.player)
	if v1 >= 4 {
		return math.MaxInt32
	}
	v2 = board.checkAlign(pos, -board.player)
	if v2 >= 4 {
		return math.MinInt32
	}
	/*
	v3 = board.checkCaptures(pos, board.player)
	if v3 > 0 {
		// TODO: Refacto this kind of things (board.capturesNb[board.player + 1], Yuck!)
		if board.capturesNb[board.player + 1] + v3 >= 10 {
			return math.MaxInt32
		}
		return board.capturesNb[board.player + 1] + v3 + 2
	}
	if (v2 * 2) > v1 {
		return v2
	}
	return v1
	*/
	return v1 + v2 * 2 + board.capturesNb[board.player + 1] * 2 - board.capturesNb[-board.player + 1] * 2
}

func (board *AIBoard) checkCaptures(pos Position, player int) int {

	defer timeFunc(time.Now(), "checkCaptures")

	x, y := pos.x, pos.y
	capt := func (incx, incy int) int {
		if !isInBounds(x + 3 * incx, y + 3 * incy) {
			return 0
		}
		if	board.board[y + incy][x + incx] == -player &&
			board.board[y + 2 * incy][x + 2 * incx] == -player &&
		 	board.board[y + 3 * incy][x + 3 * incx] == player {
				return 2
		}
		return 0
	}
	return  capt(-1, -1) + capt(1, 1) + capt(1, -1) + capt(-1, 1) +
			capt(0, -1) + capt(0, 1) + capt(-1, 0) + capt(1, 0)
}

func (board *AIBoard) checkAlign(pos Position, player int) int {
	defer timeFunc(time.Now(), "checkAlign")

	f := func (incx, incy, x, y int) int {
		cnt := 0
		x, y = x + incx, y + incy
		for i := 0; i < 4; i++ {
			if !isInBounds(x, y) || board.board[y][x] == -player {
				return cnt
			}
			if board.board[y][x] == player {
				cnt += 1
			}
			x += incx
			y += incy
		}
		return cnt
	}
	max, t := 0, 0
	x, y := pos.x, pos.y

	max = f(-1, -1, x, y) + f(1, 1, x, y)
	t = f(1, -1, x, y) + f(-1, 1, x, y)
	if t > max {
		max = t
	}	
	t = f(0, -1, x, y) + f(0, 1, x, y)
	if t > max {
		max = t
	}
	t = f(-1, 0, x, y) + f(1, 0, x, y)
	if t > max {
		return t
	} 
	return max
}
