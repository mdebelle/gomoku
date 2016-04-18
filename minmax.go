package main

import "fmt"

func search(values *[19][19]int, player, x, y, depth int, capture *[3]int) (int, int, [19][19][3]int) {

	var	score int
	var	score_a, ax, ay int
	var score_b, bx, by int

	score_a, score_b = 0, 0

	copy := [19][19][3]int {{{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0}},
							{{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0}},
							{{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0}},
							{{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0}},
							{{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0}},
							{{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0}},
							{{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0}},
							{{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0}},
							{{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0}},
							{{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0}},
							{{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0}},
							{{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0}},
							{{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0}},
							{{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0}},
							{{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0}},
							{{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0}},
							{{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0}},
							{{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0}},
							{{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0},{0,0,0}}}

	for incy := y - 5; incy < y + 6; incy++ {
//	for incy := 0; incy < 19; incy++ {
		if incy < 0 { incy = 0 } else if incy > 18 { break }
		for incx := x - 4; incx < x + 5; incx++ {
//		for incx := 0; incx < 19; incx++ {
			if incx < 0 { incx = 0 } else if incx > 18 { break }
			if values[incy][incx] == 0 {
				score = evaluateBoard(values, incx, incy, player, &copy, capture)
				fmt.Printf("score [%d][%d] -> %d\n", incx, incy, score)
				if score >= 10 {
					return incx, incy, copy
				}
				if score > score_a {
					score_a = score
					ax = incx
					ay = incy
				} else if score < score_b {
					score_b = score
					bx = incx
					by = incy
				}
			}
		}
	}

	if score_a > -score_b {
		return ax, ay, copy
	}
	return bx, by, copy
}

func checkAlign(values *[19][19]int, x, y, player int) int {
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

func checkCapt(values *[19][19]int, x, y, player int) int {
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


func evaluateBoard(values *[19][19]int, x, y, player int, copy *[19][19][3]int, capture *[3]int) int {
	
	var v1, v2, v3 int

	v1 = checkAlign(values, x, y, player)
	copy[y][x][0] = v1
	fmt.Printf("val1 -> %d\n", v1)
	if v1 >= 4 {
		return 10
	}
	v2 = -checkAlign(values, x, y, -player )
	copy[y][x][1] = -v2
	fmt.Printf("val2 -> %d\n", v2)
	if v2 <= -5 {
		return -10
	}
	v3 = checkCapt(values, x, y, player)
	copy[y][x][2] = v3
	fmt.Printf("val3 -> %d\n", v3)
	if v3 > 0 {
		return capture[player + 1] + v3
	}
	if (-v2 * 2) > v1 {
		return v2
	}
	return v1
}

	// max := 0
	// copy := *values
	// var x, y int

	// for i := 0; i < 19; i++ {
	// 	for j := 0; j < 19; j++ {
	// 		if copy[i][j] == 0 {
	// 			copy[i][j] = player
	// 			pts := 0
	// 			if checkVictory(&copy, player, i, j) {
	// 				return j, i
	// 			} else if doCaptures(&copy, player, i, j) > 0 {
	// 				pts = capturedByIA
	// 			}
	// 			p := minimise(&copy, -player)
	// 			if (pts < p) {
	// 				pts = p
	// 			}
	// 			if max < pts {
	// 				fmt.Printf("coordonees [%d][%d] = %d \n", j, i, pts)
	// 				max = pts
	// 				x = j
	// 				y = i
	// 			}
	// 			copy = *values
	// 		}
	// 	}
	// }
	// return x, y

// func minimise(values *[19][19]int, player int) int {
// 	min := 20
// 	copy := *values

// 	for i := 0; i < 19; i++ {
// 		for j := 0; j < 19; j++ {
// 			if copy[i][j] == 0 {
// 				copy[i][j] = player
// 				pts := 20
// 				if checkVictory(&copy, player, i, j) {
// 					return winPlayer
// 				} else if doCaptures(&copy, player, i, j) > 0 {
// 					pts = capturedByPlayer
// 				}
// 				p := maximise(&copy, -player)
// 				if (pts > p) {
// 					pts = p
// 				}
// 				if min > pts {
// 					min = pts
// 				}
// 				copy = *values
// 			}
// 		}
// 	}
// 	if min == 20 {
// 		return nothing
// 	}
// 	return min
// }

// func maximise(values *[19][19]int, player int) int {
// 	max := 0
// 	copy := *values

// 	for i := 0; i < 19; i++ {
// 		for j := 0; j < 19; j++ {
// 			if copy[i][j] == 0 {
// 				copy[i][j] = player
// 				pts := 0
// 				if checkVictory(&copy, player, i, j) {
// 					fmt.Printf("%v\n%d %d //%d\n", copy, j, i, player)
// 					return winIA
// 				} else if doCaptures(&copy, player, i, j) > 0{
// 					pts = capturedByIA 
// 				} else {
// 					pts = nothing
// 				}
// 				if max < pts {
// 					max = pts
// 				}
// 				copy = *values
// 			}
// 		}
// 	}
// 	if max == 0 {
// 		return nothing
// 	}
// 	return max
// }

