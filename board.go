package main

import (
	"time"
	"math"
	"sort"
	"fmt"
)

var _ = fmt.Println

type Move struct {
	pos				Position
	captures		[]Position
	forcedCaptures	[]Position
	aligns			[]AlignScore
	score			int
	isWin			bool
	isForced		bool
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

func (this *Move) IsWin() bool {
	return this.isWin
}

func (this *Move) Score() int {
	return this.score
}

func (this *Move) Position() Position {
	return this.pos
}

func (this *Move) Evaluate(board *AIBoard) {

	if debug { defer timeFunc(time.Now(), "evaluateMove") }

	// TODO: Do captures before evaluating

	if board.CanWin(this.pos) {
		alignType, forcedCaptures := checkVictory(&board.board, &board.capturesNb, this.pos.x, this.pos.y, board.player)
		if alignType == winningAlignment {
			this.score = math.MaxInt32
			this.isWin = true
			return
		} else if alignType == capturableAlignment {
			this.forcedCaptures = forcedCaptures
			this.score = math.MaxInt32
			return
		}
	}

	//p_tmp := board.board[this.pos.y][this.pos.x]
	//if p_tmp != 0 { clearAlign(&board.board, &board.alignTable, []Position{this.pos}, p_tmp) }

	_, a2, _, _ := getBestScore(&board.board, &board.alignTable, this.pos.x, this.pos.y, board.player)
	v1 := 100 / a2
	_, b2, _, _ := getBestScore(&board.board, &board.alignTable, this.pos.x, this.pos.y, -board.player)
	v2 := 100 / b2

	//if p_tmp != 0 { updateAlign(&board.board, &board.alignTable, this.pos.x, this.pos.y, p_tmp) }

	if board.capturesNb[board.player + 1] + len(this.captures) >= 10 {
		this.score = math.MaxInt32
		this.isWin = true
		return
	}

	//v1 = board.checkAlign(this.pos, board.player)
	//v2 = board.checkAlign(this.pos, -board.player)

	myCaptNb := (board.capturesNb[board.player + 1] + len(this.captures))
	hisCaptNb :=  board.capturesNb[-board.player + 1]

	this.score = v1 + v2 + myCaptNb * 15 - hisCaptNb * 15
}

type AIBoard struct {
	board		Board
	freeThrees	[2]Board
	alignTable	[2]Board
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

func (board *AIBoard) MyCapturesNb() int {
	return board.capturesNb[board.player + 1]
}

func (board *AIBoard) HisCapturesNb() int {
	return board.capturesNb[-board.player + 1]
}

func (board *AIBoard) Player() *int {
	return &board.player
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

func (this *AIBoard) GetNextMoves(forcedCaptures []Position) []Move {
	var positions []Position
	isForced := forcedCaptures != nil
	if isForced {
		positions = forcedCaptures
	} else {
		positions = this.GetSearchSpace()
	}
	moves := make([]Move, 0, len(positions))
	for _, pos := range(positions) {
		moves = append(moves, this.CreateMove(pos, isForced))
	}
	sort.Sort(sort.Reverse(ByScore(moves)))
	return moves
}

func (this *AIBoard) CreateMove(pos Position, isForced bool) Move {
	captures := make([]Position, 0, 16)
	getCaptures(&this.board, pos.x, pos.y, this.player, &captures)
	move := Move{pos, captures, nil, []AlignScore{}, 0, false, isForced}
	move.Evaluate(this)
	return move
}

func (this *AIBoard) DoMove(move Move) {
	this.board[move.pos.y][move.pos.x] = this.player
	doCaptures(&this.board, &move.captures)
	this.capturesNb[this.player + 1] += len(move.captures)
	clearAlign(&this.board, &this.alignTable, move.captures, -this.player)
	updateAlign(&this.board, &this.alignTable, move.pos.x, move.pos.y, this.player)
}

func (this *AIBoard) UndoMove(move Move) {
	this.board[move.pos.y][move.pos.x] = empty
	undoCaptures(&this.board, &move.captures, this.player)
	this.capturesNb[this.player + 1] -= len(move.captures)
	clearAlign(&this.board, &this.alignTable, []Position{move.pos}, this.player)
	for _, pos := range move.captures {
		updateAlign(&this.board, &this.alignTable, pos.x, pos.y, -this.player)
	}
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

func NewAIBoard(values *Board, freeThree, alignTable *[2]Board, capture *[3]int, player, depth int) AIBoard {
	board := AIBoard{*values, *freeThree, *alignTable, *capture, player, depth}
	return board
}

func (this *AIBoard) CanWin(pos Position) bool {
	// TODO: Redo this
	a1, a2, a3, a4 := getScore(&this.alignTable, pos.x, pos.y, this.player)
	for _, p := range []int{a1, a2, a3, a4} {
		if p == 15 || p == 22 || p == 24 {
			return true
		}
	}
	return false
}

func (board *AIBoard) Evaluate(move *Move) (score int, quiet bool) {

	if debug { defer timeFunc(time.Now(), "evaluateBoard") }

	// WIP: Quiescence
	// Also check if there is forced captures on this turn
	// Maybe put a bolean somewhere

	// TODO: Already done
	if board.CanWin(move.pos) {
		alignType, _ := checkVictory(&board.board, &board.capturesNb, move.pos.x, move.pos.y, board.player)
		if alignType == winningAlignment {
			return math.MaxInt32 * 2, true
		} else if alignType == capturableAlignment {
			return math.MaxInt32, false
		}
	}

	if board.capturesNb[board.player + 1] >= 10 {
		return math.MaxInt32, true
	}

	quiet = move.forcedCaptures == nil

	eval := func (x, y int) int {
		_, a2, _, _ := getBestScore(&board.board, &board.alignTable, x, y, board.player)
		v1 := 100 / a2
		_, b2, _, _ := getBestScore(&board.board, &board.alignTable, x, y, -board.player)
		v2 := 100 / b2

		return v1 - v2
	}

	score = 0

	// TODO: GetSearchSpace is now too low for this use
	// TODO: Maybe update total with each move
	positions := board.GetSearchSpace()
	for _, pos := range positions {
		score += eval(pos.x, pos.y)
	}

	return score + board.MyCapturesNb() * 100 - board.HisCapturesNb() * 100, quiet
}

/*
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
*/
