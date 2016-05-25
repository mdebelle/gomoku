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

	alpha = -math.MaxInt32
	beta = math.MaxInt32

	s := -searchdeeper(b, move, 1, alpha, beta)
	/*
	fmt.Println("-----------------")
	fmt.Println(b.board)
	fmt.Println("Score: ", move.Score())
	fmt.Printf("||move[%d][%d], %d\n", move.pos.y, move.pos.x, s)
	//*/
	c <- s
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
		boardData[move.pos.y][move.pos.x][6] = 1
		boardData[move.pos.y][move.pos.x][7] = move.score
		if move.IsWin() {
			return move.pos.x, move.pos.y, boardData
		}

		b.DoMove(move)
		b.UpdateFreeThrees(move.pos, move.captures)
		// TODO: Multithreading
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
	} else if depth == 2 {
		return bestLeaf(b, move, depth, alpha, beta)
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

func bestLeaf(b *AIBoard, move *Move, depth, alpha, beta int) int {

	bestscore := math.MinInt32

	b.SwitchPlayer()
	defer b.SwitchPlayer()

	moves := b.GetNextMoves(move.forcedCaptures)

	var chans []chan int

	for i := 0; i < len(moves); i++ {
		chans = append(chans, make(chan int))
	}

	for i, move := range(moves) {
		if move.IsWin() {
			return move.Score()
		}
		b.DoMove(move)

		var move2 = move
		var b2 = *b
		go multisearch(&b2, &move2, -beta, -alpha, chans[i])
		b.UndoMove(move)
	}

	for i:= 0; i < len(moves); i++ {
		s := <-chans[i]
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

func searchdeeper2(b *AIBoard, move *Move, depth, alpha, beta int) int {

	nodesSearched++
	bestscore := math.MinInt32

	b.depth++
	defer func() {b.depth--}()

	b.SwitchPlayer()
	defer b.SwitchPlayer()

	moves := b.GetNextMoves(move.forcedCaptures)

	for _, move := range(moves) {
		if move.IsWin() {
			return move.Score()
		}
		b.DoMove(move)
//		s := -move.Score()
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
