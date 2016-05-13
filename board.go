package main

import (
	"time"
	"math"
	"sort"
)

type Move struct {
	pos			Position
	captures	[]Position
	score		int
}

type ByScore []Move

func (this ByScore) Len() int {
	return len(this)
}

func (this ByScore) Swap(a, b int) {
	this[a], this[b] = this[b], this[a]
}

func (this ByScore) Less(a, b int) bool {
	return this[a].score < this[b].score
}

func (this *Move) Score() int {
	return this.score
}

func (this *Move) Position() Position {
	return this.pos
}

type AIBoard struct {
	board		Board
	//alignments	[2][4]Board
	freeThrees	[2]Board
	capturesNb	[3]int
	player		int
	depth		int
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

	if debug { defer timeFunc(time.Now(), "getSearchSpace") }

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
				for i := 1; i < 2; i++ {
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

func (this *AIBoard) GetNextMoves() []Move {
	positions := this.GetSearchSpace()
	moves := make([]Move, 0, len(positions))
	for _, pos := range(positions) {
		moves = append(moves, this.CreateMove(pos))
	}
	sort.Sort(sort.Reverse(ByScore(moves)))
	return moves
}

func (this *AIBoard) CreateMove(pos Position) Move {
	captures := make([]Position, 0, 16)
	getCaptures(&this.board, pos.x, pos.y, this.player, &captures)
	move := Move{pos, captures, 0}
	this.DoMove(move)
	score := this.Evaluate(pos)
	this.UndoMove(move)
	move.score = score
	return move
}

func (this *AIBoard) DoMove(move Move) {
	this.board[move.pos.y][move.pos.x] = this.player
	doCaptures(&this.board, &move.captures)
	this.capturesNb[this.player + 1] += len(move.captures)
}

func (this *AIBoard) UndoMove(move Move) {
	this.board[move.pos.y][move.pos.x] = empty
	undoCaptures(&this.board, &move.captures, this.player)
	this.capturesNb[this.player + 1] -= len(move.captures)
}

func (board *AIBoard) UpdateFreeThrees(pos Position, captures []Position) {
	// TODO: Two functions -> update from move and update from move cancelation
	// TODO: And store changes. In fact, find a new method for free threes calculation / examination
	if debug { defer timeFunc(time.Now(), "updateFreeThree") }

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

	if debug { defer timeFunc(time.Now(), "evaluateBoard") }

	var v1, v2 int

	v1 = board.checkAlign(pos, board.player)
	if v1 >= 4 || board.capturesNb[board.player + 1] >= 10 {
		return math.MaxInt32 + board.depth
	}

	v2 = board.checkAlign(pos, -board.player)

	return v1 + v2 * 2 + board.capturesNb[board.player + 1] * 2 - board.capturesNb[-board.player + 1] * 2
}

func (board *AIBoard) checkCaptures(pos Position, player int) int {

	if debug { defer timeFunc(time.Now(), "checkCaptures") }

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
	if debug { defer timeFunc(time.Now(), "checkAlign") }

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
