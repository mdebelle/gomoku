package main

import (
	"fmt"
	"time"
	"math"
)

var nodesSearched = 0

func printAlignments(b *AIBoard) {
	for y := 0; y < 19; y++ {
		for x := 0; x < 19; x++ {
			a, e, c, d := getScore(&b.alignTable, x, y, player_one)
			fmt.Printf("[%v ", a + e + c + d)
			a, e, c, d = getScore(&b.alignTable, x, y, player_two)
			fmt.Printf("%v]", a + e + c + d)
		}
		fmt.Println()
	}
}

func search(values *Board, freeThree, alignTable *[2]Board, player, x, y, depth int, capture *[3]int, forcedCaptures []Position) (int, int, BoardData) {

	nodesSearched = 1
	startTime := time.Now()
	defer func () {fmt.Println(nodesSearched, "nodes searched in", time.Since(startTime), "(", time.Since(startTime) / time.Duration(nodesSearched), "by node)")}()

	var	ax, ay int
	var boardData BoardData

	b := NewAIBoard(values, freeThree, alignTable, capture, player, depth)

	moves := b.GetNextMoves(forcedCaptures)

	ax, ay = moves[0].pos.x, moves[0].pos.y

	bestscore := math.MinInt32
	alpha := math.MinInt32
	beta := math.MaxInt32

	for _, move := range(moves) {
		boardData[move.pos.y][move.pos.x][6] = 1
		boardData[move.pos.y][move.pos.x][7] = move.score
		if move.IsWin() {
			return move.pos.x, move.pos.y, boardData
		}

		b.DoMove(move)
		b.UpdateFreeThrees(move.pos, move.captures)
		// TODO: Multithreading
		s := -searchdeeper(&b, &move, -beta, -alpha)
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

func searchdeeper(b *AIBoard, move *Move, alpha, beta int) int {

	nodesSearched++
	bestscore := math.MinInt32

	b.depth--
	defer func() {b.depth++}()

	if b.depth == 0 {
		score, quiet := b.Evaluate(move)
		//score, quiet := move.Score(), true
		if (!quiet) {
			b.depth += 3
			score = searchdeeper(b, move, alpha, beta)
			//fmt.Println("GROSSE BITE", score)
			b.depth -= 3
			return score
		}
		return -score
	}

	b.SwitchPlayer()
	defer b.SwitchPlayer()

	moves := b.GetNextMoves(move.forcedCaptures)

	for _, move := range(moves) {
		if move.IsWin() {
			return move.Score()
		}
		b.DoMove(move)
		s := -searchdeeper(b, &move, -beta, -alpha)
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
