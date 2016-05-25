package main

import (
	"fmt"
	"time"
	"math"
)

var nodesSearched = 0

var timer time.Time

func search(values *Board, freeThree, alignTable *[2]Board, player, x, y, depth int, capture *[3]int, forcedCaptures []Position) (int, int, BoardData) {

	timer = time.Now()

	nodesSearched = 1
	startTime := time.Now()
	defer func () {fmt.Println(nodesSearched, "nodes searched in", time.Since(startTime), "(", time.Since(startTime) / time.Duration(nodesSearched), "by node)")}()

	var	ax, ay int
	var boardData BoardData

	b := NewAIBoard(values, freeThree, alignTable, capture, player, 0)

	moves := b.GetNextMoves(forcedCaptures)

	ax, ay = moves[0].pos.x, moves[0].pos.y

	bestscore := math.MinInt32
	alpha := math.MinInt32
	beta := math.MaxInt32

	for _, move := range(moves) {
//		fmt.Println(move)
		boardData[move.pos.y][move.pos.x][6] = 1
		boardData[move.pos.y][move.pos.x][7] = move.score
		if move.IsWin() {
			return move.pos.x, move.pos.y, boardData
		}

		b.DoMove(move)
		b.UpdateFreeThrees(move.pos, move.captures)
		s := -searchdeeper(&b, &move, depth - 1, -beta, -alpha)
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

	return ax, ay, boardData
}

func searchdeeper(b *AIBoard, move *Move, depth, alpha, beta int) int {

	nodesSearched++
	bestscore := math.MinInt32

	b.depth++
	defer func() {b.depth--}()

	if depth == 0 {
		if move.isForced {
			return math.MaxInt32
		}
		return -move.Score()
	} else if time.Since(timer) >= time.Millisecond * 498 {
		return 0
		//return -move.Score()
	}

	b.SwitchPlayer()
	defer b.SwitchPlayer()

	moves := b.GetNextMoves(move.forcedCaptures)

	for i, move := range(moves) {
		if i > 12 {
			break
		}
		if move.IsWin() {
			return move.Score()
		}
		b.DoMove(move)
		s := -searchdeeper(b, &move, depth - 1, -beta, -alpha)
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
