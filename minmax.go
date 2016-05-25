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

func multisearch(b *AIBoard, pos *Position, c chan Move) {
	m := b.CreateMove(*pos, false)
	c <- m
}

func search(values *Board, freeThree, alignTable *[2]Board, player, x, y, depth int, capture *[3]int, forcedCaptures []Position) (int, int, BoardData) {

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
		//*
		//move.Evaluate(b)
		score := move.Score()
		return -score
		/*/
		score, quiet := b.Evaluate(move)
		if (!quiet) {
			depth += 2
		} else {
			return -score
		}
		//*/

	} /*else if depth == 1 {
		return bestLeaf(b, move, depth - 1, -beta, -alpha)
	} //*/

	b.SwitchPlayer()
	defer b.SwitchPlayer()

	moves := b.GetNextMoves(move.forcedCaptures)

	for _, move := range(moves) {
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

func bestLeaf(b *AIBoard, move *Move, depth, alpha, beta int) int {

	bestscore := math.MinInt32

	b.SwitchPlayer()
	defer b.SwitchPlayer()

	var movesPositions []Position
	if (move.forcedCaptures != nil) {
		movesPositions = move.forcedCaptures
	} else {
		movesPositions = b.GetSearchSpace()
	}

	var chans []chan Move

	for i := 0; i < len(movesPositions); i++ {
		chans = append(chans, make(chan Move))
	}

	for i, pos := range(movesPositions) {
		go multisearch(b, &pos, chans[i])
	}

	for i:= 0; i < len(movesPositions); i++ {
		m := <-chans[i]
		if m.IsWin() {
			return m.Score()
		}

		s := m.Score()

		//*
		if s >= beta {
			return s
		}
		//*/

		if s > bestscore {
			bestscore = s
			if s > alpha {
				alpha = s
			}
		}
	}

	return bestscore
}
