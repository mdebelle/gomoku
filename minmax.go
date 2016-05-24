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

func multisearch(b *AIBoard, move *Move, alpha, beta int, c chan int) {

	s := -searchdeeper(b, move, 5 - 1, alpha, beta)
	c <- s	
//	fmt.Printf("||move[%d][%d], %d\n", move.pos.y, move.pos.x, s)
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

	var chans []chan int

	for i := 0; i < len(moves); i++ {
		chans = append(chans, make(chan int))
	}

	for i, move := range(moves) {
		boardData[move.pos.y][move.pos.x][6] = 1
		boardData[move.pos.y][move.pos.x][7] = move.score
		if move.IsWin() {
			return move.pos.x, move.pos.y, boardData
		}

		b.DoMove(move)
		b.UpdateFreeThrees(move.pos, move.captures)
		//fmt.Printf("%d//%v\n", i, move)
		var move2 = move
		var b2 = b

		go multisearch(&b2, &move2, -beta, -alpha, chans[i])
//		s := -searchdeeper(&b, &move, -beta, -alpha)
//		boardData[move.pos.y][move.pos.x][5] = s
		b.UndoMove(move)
		b.UpdateFreeThrees(move.pos, move.captures)

		// if s >= beta {
		// 	return move.pos.x, move.pos.y, boardData
		// }

		// if s > bestscore {
		// 	bestscore = s
		// 	ax, ay = move.pos.x, move.pos.y
		// 	if s > alpha {
		// 		alpha = s
		// 	}
		// }
	}

	for i:= 0; i < len(moves); i++ {
		s := <-chans[i]
		//fmt.Printf("move[%d][%d], %d\n", moves[i].pos.y, moves[i].pos.x, s)
		boardData[moves[i].pos.y][moves[i].pos.x][5] = s
		if s >= beta {
			return moves[i].pos.x, moves[i].pos.y, boardData
		}

		if s > bestscore {
			bestscore = s
			ax, ay = moves[i].pos.x, moves[i].pos.y
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
		//score, quiet := b.Evaluate(move)
		score, quiet := move.Score(), true
		if (!quiet) {
			//fmt.Println("Not quiet", b.depth)
			//depth += 2
			//score = -searchdeeper(b, move, 1, alpha, beta)
			//fmt.Println(-score)
			return -score
		} else {
			return -score
		}
	}

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
