package main

//import "fmt"
//import "math"

type Position struct {
	x, y int
}

func search(values *[19][19]int, player, x, y, depth int, capture *[3]int) (int, int, [19][19][3]int) {

	var	score_a, ax, ay int
	var	score_b, bx, by int
	var pos *Position
	lst := []*Position{}


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

	var copy_capt [19][19]int

	score_a, score_b = 0, 0

	for circle := 1; circle < 7; circle++ {
		a, b := y - circle, y + circle
		for incx := x - (circle - 1); incx < x + circle; incx++ {
			if incx < 0 { incx = 0 } else if incx > 18 { break }
			if a >= 0 && values[a][incx] == 0 {
				pos = new(Position)
				pos.x, pos.y = incx, a
				lst = append(lst, pos)
			}
			if b < 19 && values[b][incx] == 0 {
				pos = new(Position)
				pos.x, pos.y = incx, b
				lst = append(lst, pos)
			}
		}
		a, b = x - circle, x + circle
		for incy := y - circle; incy <= y + circle; incy++ {
			if incy < 0 { incy = 0 } else if incy > 18 { break }
			if a >= 0 && values[incy][a] == 0 {
				pos = new(Position)
				pos.x, pos.y = a, incy
				lst = append(lst, pos)
			}
			if b < 19 && values[incy][b] == 0 {
				pos = new(Position)
				pos.x, pos.y = b, incy
				lst = append(lst, pos)
			}
		}
	}

	for i := range(lst) {
		score := evaluateBoard(values, lst[i].x, lst[i].y, player, &copy, capture)
		if score >= 20 {
			return lst[i].x, lst[i].y, copy
		}
		
		do_move(lst[i].x, lst[i].y, copy[lst[i].y][lst[i].x][2], values, &copy_capt, player)
//		searchdeeper(values, -player, x, y, depth - 1, capture)
		undo_move(lst[i].x, lst[i].y, copy[lst[i].y][lst[i].x][2], values, &copy_capt)
			
		
		if score > score_a {
			score_a, ax, ay = score, lst[i].x, lst[i].y
		} else if score < score_b {
			score_b, bx, by = score, lst[i].x, lst[i].y
		}	
	}
	if score_a > -score_b {
		return ax, ay, copy
	}
	return bx, by, copy
}

func searchdeeper(values *[19][19]int, player, x, y, depth int, capture *[3]int) int {

	var	score_a int
	var	score_b int
	var pos *Position
	lst := []*Position{}

//	var copy_capt [19][19]int
	var copy [19][19][3]int

	score_a, score_b = 0, 0

	for circle := 1; circle < 7; circle++ {
		a, b := y - circle, y + circle
		for incx := x - (circle - 1); incx < x + circle; incx++ {
			if incx < 0 { incx = 0 } else if incx > 18 { break }
			if a >= 0 && values[a][incx] == 0 {
				pos = new(Position)
				pos.x, pos.y = incx, a
				lst = append(lst, pos)
			}
			if b < 19 && values[b][incx] == 0 {
				pos = new(Position)
				pos.x, pos.y = incx, b
				lst = append(lst, pos)
			}
		}
		a, b = x - circle, x + circle
		for incy := y - circle; incy <= y + circle; incy++ {
			if incy < 0 { incy = 0 } else if incy > 18 { break }
			if a >= 0 && values[incy][a] == 0 {
				pos = new(Position)
				pos.x, pos.y = a, incy
				lst = append(lst, pos)
			}
			if b < 19 && values[incy][b] == 0 {
				pos = new(Position)
				pos.x, pos.y = b, incy
				lst = append(lst, pos)
			}
		}
	}

	for i := range(lst) {
		score := evaluateBoard(values, lst[i].x, lst[i].y, player, &copy, capture)
		if score >= 20 {
			return score
		}

		// if depth > 0 && ((depth % 2 == 1 && (score > score_a || score < score_b)) || (depth % 2 == 0 && (score < score_a && score > score_b))) {
		// 	a := do_move(lst[i].x, lst[i].y, copy[lst[i].y][lst[i].x][2], values, &copy_capt, player)
		// 	capture[player + 1] += a
		// 	searchdeeper(values, -player, x, y, depth - 1, capture)
		// 	undo_move(lst[i].x, lst[i].y, copy[lst[i].y][lst[i].x][2], values, &copy_capt)
		// 	capture[player + 1] -= a
		// }

		if score > score_a {
			score_a = score
		} else if score < score_b {
			score_b = score
		}
	}
	if score_a > -score_b {
		return score_a
	}
	return score_b
}

func do_move(x, y, capt int, values *[19][19]int, copy_capt *[19][19]int, player int) {

	if (capt > 1) {
		copy_capt = values
		doCaptures(values, player, y, x)
	} else {
		values[y][x] = player
	}
}

func undo_move(x, y, capt int, values *[19][19]int, copy_capt *[19][19]int) {

	if (capt > 1) {
		values = copy_capt
	} else {
		values[y][x] = 0
	}

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
		return capture[player + 1] + v3 + 2
	}
	if (-v2 * 2) > v1 {
		return v2
	}
	return v1
}
