package main

//import "math"
//import "time"

import "fmt"

type Position struct {
	x, y int
}

// var stopByTime = false
// var node = 0

func search(values, freeThree *Board, player, x, y, depth int, capture *[3]int) (int, int, BoardData) {

	var	ax, ay int
//	var pos *Position
//	lst := []*Position{}
	var copy BoardData
	var copy_capt Board

	checkfree := func (b int) bool {
		var j int
		for i := uint(0); i < 4; i++ {
			if b & (1 << i) != 0 { j++}
		}
		if j == 2 {return true}
		return false
	}
	
	m := make(map[Position]bool)

	for i:= 0;i < 19; i++ {
		for j:= 0;j < 19; j++ {
			if values[i][j] != 0 {

				for circle := 1; circle < 3; circle++ {
					a, b := i - circle, i + circle
					for incx := j - (circle - 1); incx < j + circle; incx++ {
						if incx < 0 { incx = 0 } else if incx > 18 { break }
						if a >= 0 && values[a][incx] == 0 && !checkfree(freeThree[a][incx]) {
							m[Position{incx, a}] = true
						}
						if b < 19 && values[b][incx] == 0 && !checkfree(freeThree[b][incx]) {
							m[Position{incx, b}] = true
						}
					}
					a, b = j - circle, j + circle
					for incy := i - circle; incy <= i + circle; incy++ {
						if incy < 0 { incy = 0 } else if incy > 18 { break }
						if a >= 0 && values[incy][a] == 0 && !checkfree(freeThree[incy][a]) {
							m[Position{a, incy}] = true
						}
						if b < 19 && values[incy][b] == 0 && !checkfree(freeThree[incy][b]) {
							m[Position{b, incy}] = true
						}
					}
				}

			}
		}
	}

	fmt.Printf("%v\n", m)

	bestscore := -int(^uint32(0)>>1) // int le plus large possible dans les negatifs

	alpha := -int(^uint32(0)>>1)
	beta := int(^uint32(0)>>1)
	fmt.Printf("-> %d | %d\n", alpha, beta)

	for i, _ := range(m) {
		s := evaluateBoard(values, i.x, i.y, player, &copy, capture)
		if s >= 20 {
			return i.x, i.y, copy
		}
		a := do_move(i.x, i.y, copy[i.y][i.x], values, &copy_capt, player)
		capture[player+1] += a
		s = -searchdeeper(values, -player, x, y, depth - 1, capture, -beta, -alpha)
		undo_move(i.x, i.y, copy[i.y][i.x], values, &copy_capt)
		capture[player+1] -= a
		
		if s >= beta {
			return i.x, i.y, copy
		}
		if s > bestscore {
			bestscore = s
			ax, ay = i.x, i.y
			if s > alpha {
				alpha = s
			}
		}
	}

	return ax, ay, copy

}

func searchdeeper(values *Board, player, x, y, depth int, capture *[3]int, alpha, beta int) int {

//	var pos *Position
//	lst := []*Position{}
	var copy_capt Board
	var copy BoardData

	if depth == 0 {
		return evaluateBoard(values, x, y, player, &copy, capture)
	}

	m := make(map[Position]bool)

	for i:= 0;i < 19; i++ {
		for j:= 0;j < 19; j++ {
			
			if values[i][j] != 0 {

				for circle := 1; circle < 3; circle++ {
					a, b := i - circle, i + circle
					for incx := j - (circle - 1); incx < j + circle; incx++ {
						if incx < 0 { incx = 0 } else if incx > 18 { break }
						if a >= 0 && values[a][incx] == 0 {
							m[Position{incx, a}] = true
						}
						if b < 19 && values[b][incx] == 0 {
							m[Position{incx, b}] = true
						}
					}
					a, b = j - circle, j + circle
					for incy := i - circle; incy <= i + circle; incy++ {
						if incy < 0 { incy = 0 } else if incy > 18 { break }
						if a >= 0 && values[incy][a] == 0 {
							m[Position{a, incy}] = true
						}
						if b < 19 && values[incy][b] == 0 {
							m[Position{b, incy}] = true
						}
					}
				}

			}
		}
	}

	bestscore := -int(^uint32(0)>>1) // int le plus large possible dans les negatif
	// node++

	// if node % 4095 == 0 {
	// 	if time.Now().After(stopTime) {
	// 		stopByTime = true
	// 		return 0
	// 	}
	// }


	for i, _ := range(m) {
		score := evaluateBoard(values, i.x, i.y, player, &copy, capture)
		if score >= 20 {
			return score
		}
	//	if stopByTime { break }
		a := do_move(i.x, i.y, copy[i.y][i.x], values, &copy_capt, player)
		capture[player+1] += a
		s := -searchdeeper(values, -player, x, y, depth - 1, capture, -beta, -alpha)
		undo_move(i.x, i.y, copy[i.y][i.x], values, &copy_capt)
		capture[player+1] -= a
		
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

func do_move(x, y int, capt [5]int, values *Board, copy_capt *Board, player int) int {

	if (capt[2] > 1) {
		*copy_capt = *values
		return doCaptures(values, player, y, x)
	} else {
		values[y][x] = player
	}
	return 0
}

func undo_move(x, y int, capt [5]int, values *Board, copy_capt *Board) {

	if (capt[2] > 1) {
		*values = *copy_capt
	} else {
		values[y][x] = 0
	}

}

func checkAlign(values *Board, x, y, player int) int {
	f := func (incx, incy, x, y, p int) int {
		cnt := 0
		x, y = x + incx, y + incy
		for i := 0; i < 4; i++ {
			if !checkBounds(x, y) || values[y][x] == -p{
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
		if !checkBounds(x + 3 * incx, y + 3 * incy) {
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

func evaluateBoard(values *Board, x, y, player int, copy *BoardData, capture *[3]int) int {

	var v1, v2, v3 int

	v1 = checkAlign(values, x, y, player)
	copy[y][x][0] = v1
	if v1 >= 4 {
		return 20
	}
	v2 = -checkAlign(values, x, y, -player )
	copy[y][x][1] = -v2
	if v2 <= -4 {
		return -20
	}
	v3 = checkCapt(values, x, y, player)
	copy[y][x][2] = v3
	if v3 > 0 {
		if capture[player + 1] + v3 >= 10 {
			return 20
		}
		return capture[player + 1] + v3 + 2
	}
	if (-v2 * 2) > v1 {
		return v2
	}
	return v1
}
