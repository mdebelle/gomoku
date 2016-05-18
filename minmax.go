package main

import (
	"fmt"
	"time"
	"math"
)

var nodesSearched = 0

func search(values *Board, freeThree, alignTable *[2]Board, player, x, y, depth int, capture *[3]int) (int, int, BoardData) {

	nodesSearched = 0
	startTime := time.Now()

	var	ax, ay int
	var boardData BoardData

	b := NewAIBoard(values, freeThree, alignTable, capture, player, depth)

	moves := b.GetNextMoves()

	ax, ay = moves[0].pos.x, moves[0].pos.y

	bestscore := math.MinInt32
	alpha := math.MinInt32
	beta := math.MaxInt32

	for _, move := range(moves) {
		boardData[move.pos.y][move.pos.x][6] = 1
		boardData[move.pos.y][move.pos.x][7] = move.score
		if move.score >= 2e9 {
			return move.pos.x, move.pos.y, boardData
		}

		fmt.Println("-----------BEFORE------------")
		for y := 0; y < 19; y++ {
			for x := 0; x < 19; x++ {
				a, e, c, d := getScore(&b.alignTable, x, y, player_one)
				fmt.Printf("[%v ", a + e + c + d)
				a, e, c, d = getScore(&b.alignTable, x, y, player_two)
				fmt.Printf("%v]", a + e + c + d)
			}
			fmt.Println()
		}

		b.DoMove(move)
		b.UpdateFreeThrees(move.pos, move.captures)
		s := -searchdeeper(&b, move.pos, depth - 1, -beta, -alpha)
		boardData[move.pos.y][move.pos.x][5] = s
		b.UndoMove(move)
		b.UpdateFreeThrees(move.pos, move.captures)


		fmt.Println("-----------AFTER------------")
		for y := 0; y < 19; y++ {
			for x := 0; x < 19; x++ {
				a, e, c, d := getScore(&b.alignTable, x, y, player_one)
				fmt.Printf("[%v ", a + e + c + d)
				a, e, c, d = getScore(&b.alignTable, x, y, player_two)
				fmt.Printf("%v]", a + e + c + d)
			}
			fmt.Println()
		}

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
	return ax, ay, boardData
}

func searchdeeper(b *AIBoard, move Position, depth int, alpha, beta int) int {

	nodesSearched++
	bestscore := math.MinInt32

	b.depth--
	if depth == 0 {
		return -b.Evaluate(move)
	}

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
