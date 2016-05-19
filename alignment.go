package main

// returne weight of alignement in the selected axe
// and time to expect an wining alignment
// if space missing return -1 for time
func winingAlignement(board *Board, axe1, axe2, x, y, player int) (int, int){

	// AlignementGagnant
	if axe1 == 15 || axe2 == 15 || (axe1 == 14 && axe2 >= 8) || (axe2 == 14 && axe1 >= 8) || (axe1 >= 12 && axe2 >= 12) {
		return axe1+axe2, 0
	}

	space := 8
	chain := 0
	spaceAndChaine := func (axe int) {
		for i := uint(3); i >= 0; i-- {
			if ((axe >> i) & 1) == 0 {
				if !isInBounds(x-(4-int(i)), y) || board[y][x-(4-int(i))] == -player {
					space -= int(i+1)
					break
				}  
			} else {
				chain++
			} 
		}
	}
	spaceAndChaine(axe1)
	spaceAndChaine(axe2)

	// Possibility of wining alignment
	if chain + space >= 5 {
		if chain >= 5 {
			return axe1+axe2, 1
		} else {
			return axe1+axe2, (5 - chain)
		}
	}

	// Space missing
	return axe1+axe2, -1
}

func getLeftRightScore(board *Board, alignTable *[2]Board, x, y, player int) (int, int) {

	if !isInBounds(x,y) { return -1, -1 }

	const (
		maskleft = 0x0000000f
		maskright = 0x000000f0
	)

	v := alignTable[(player+1)/2][y][x]

	l := (v & maskleft)
	r := ((v & maskright) >> 4)

	return winingAlignement(board, l, r, x, y, player)
}

func getTopBottomScore(board *Board, alignTable *[2]Board, x, y, player int) (int, int) {

	if !isInBounds(x,y) { return -1, -1 }

	const (
		masktop = 0x00000f00
		maskbottom = 0x0000f000
	)

	v := alignTable[(player+1)/2][y][x]

	t := (v & masktop) >> 8
	b := (v & maskbottom) >> 12

	return winingAlignement(board, t, b, x, y, player)
}

func getLeftTopRightBottomScore(board *Board, alignTable *[2]Board, x, y, player int) (int, int) {

	if !isInBounds(x,y) { return -1, -1 }

	const (
		masklefttop = 0x000f0000
		maskrightbottom = 0x00f00000
	)

	v := alignTable[(player+1)/2][y][x]

	lt := ((v & masklefttop) >> 16) 
	rb := ((v & maskrightbottom) >> 20)

	return winingAlignement(board, lt, rb, x, y, player)
}

func getRightTopLeftBottomScore(board *Board, alignTable *[2]Board, x, y, player int) (int, int) {

	if !isInBounds(x,y) { return -1, -1 }

	const (
		maskrighttop = 0x0f000000
		maskleftbottom = 0xf0000000
	)

	v := alignTable[(player+1)/2][y][x]

	rt := ((v & maskrighttop) >> 24)
	lb := ((v & maskleftbottom) >> 28)

	return winingAlignement(board, rt, lb, x, y, player)
}

func getBestScore(board *Board, alignTable *[2]Board, x, y, player int) (int, int) {
	
	var s, t []int

	s[0], t[0] = getLeftRightScore(board, alignTable, x, y, player)
	s[1], t[1] = getTopBottomScore(board, alignTable, x, y, player)
	s[2], t[2] = getLeftTopRightBottomScore(board, alignTable, x, y, player)
	s[3], t[3] = getRightTopLeftBottomScore(board, alignTable, x, y, player)
	min, indexmin := 8, 0
	max, indexmax := 0, 0
	
	for i, v := range t {
		if  v < min && v >= 0 {
			min = v
			indexmin = i
		}
		if v > max {
			max = v
			indexmax = i
		}
	}
	
	return s[indexmin], s[indexmax]

}


func getScore(alignTable *[2]Board, x, y, player int) (int, int, int, int) {

	if !isInBounds(x,y) { return 0,0,0,0 }

	const (
		maskleft = 0x0000000f
		maskright = 0x000000f0
		masktop = 0x00000f00
		maskbottom = 0x0000f000
		masklefttop = 0x000f0000
		maskrightbottom = 0x00f00000
		maskrighttop = 0x0f000000
		maskleftbottom = 0xf0000000
	)

	v := alignTable[(player+1)/2][y][x]

	axeLeftRight := (v & maskleft) + ((v & maskright) >> 4)
	axeTopBottom := ((v & masktop) >> 8) + ((v & maskbottom) >> 12)
	axeLeftTopRightBottom := ((v & masklefttop) >> 16) + ((v & maskrightbottom) >> 20)
	axeRightTopLeftBottom := ((v & maskrighttop) >> 24) + ((v & maskleftbottom) >> 28)

	return axeLeftRight, axeTopBottom, axeLeftTopRightBottom, axeRightTopLeftBottom
}

func clearAlign(board *Board, alignTable *[2]Board, lst []Position, hey int) {

	clearOponentSituation := func (player, py, px, axe, start int) {
		for j := start; j < 5; j++ {
			alignTable[(-player+1)/2][py][px] &= ^(1 << uint(axe-j))
		}
	}

	resetAxeScore := func(x, y, incx, incy, axe, c int) {
		var state bool

		if !isInBounds(x,y) { return }

		player := board[y][x]

		if player != 0 {
			for i := 1; i < 5; i++ {
				if isInBounds(x+(i*incx), y+(i*incy)) {
					if !state && board[y+(i*incy)][x+(i*incx)] == 0 {
						alignTable[(player+1)/2][y+(i*incy)][x+(i*incx)] |= (1 << uint(axe-i))
						clearOponentSituation(player, y+(i*incy), x+(i*incx), axe, i)
					} else if board[y+(i*incy)][x+(i*incx)] == -player {
						state = true
					}
				}
			}
		} else {
			if axe % 8 == 0 { 
				alignTable[(hey+1)/2][y][x] &= ^(1 << uint(axe-4-c))
			} else {
				alignTable[(hey+1)/2][y][x] &= ^(1 << uint(axe+4-c))
			}
		}
	}

	const (
		axeLeft = 4
		axeRight = 8
		axeTop = 12
		axeBottom = 16
		axeLeftTop = 20
		axeRightBottom = 24
		axeRightTop = 28
		axeLeftBottom = 32
	)

	for _, p := range lst {

		alignTable[(player_one + 1)/2][p.y][p.x] = 0
		alignTable[(player_two + 1)/2][p.y][p.x] = 0

		for i := 1; i < 5; i++ {
			// LeftRight
			resetAxeScore(p.x-i, p.y, 1, 0, axeRight, i)
			resetAxeScore(p.x+i, p.y, -1, 0, axeLeft, i)
			resetAxeScore(p.x, p.y-i, 0, 1, axeBottom, i)
			resetAxeScore(p.x, p.y+i, 0, -1, axeTop, i)
			resetAxeScore(p.x-i, p.y-i, 1, 1, axeRightBottom, i)
			resetAxeScore(p.x+i, p.y+i, -1, -1, axeLeftTop, i)
			resetAxeScore(p.x+i, p.y-i, -1, 1, axeLeftBottom, i)
			resetAxeScore(p.x-i, p.y+i, 1, -1, axeRightTop, i)
		}
	}
}

func updateAlign(board *Board, alignTable *[2]Board, x, y, player int) {
	
	clearOponentSituation := func (player, py, px, axe, start int) {
		for j := start; j < 5; j++ {
			alignTable[(-player+1)/2][py][px] &= ^(1 << uint(axe-j))
		}
	}

	updateDistanceScore := func (player, px, py, axe, i int, state bool) bool {
		if isInBounds(px, py) {
			if !state && board[py][px] == 0 {
				alignTable[(player+1)/2][py][px] |= (1 << uint(axe-i))
				clearOponentSituation(player, py, px, axe, i)
			} else if board[py][px] == -player {
				state = true
			}
		}
		return state
	}

	const (
		axeLeft = 4
		axeRight = 8
		axeTop = 12
		axeBottom = 16
		axeLeftTop = 20
		axeRightBottom = 24
		axeRightTop = 28
		axeLeftBottom = 32
	)

	var l, lt, t, rt, r, rb, b, lb bool

	alignTable[(player+1)/2][y][x] = 0
	alignTable[(-player+1)/2][y][x] = 0
	
	for i:= 1; i < 5; i++ {
		// Axe left right
		l = updateDistanceScore(player, x-i, y, axeLeft, i, l)
		r = updateDistanceScore(player, x+i, y, axeRight, i, r)
		// Axe top left
		t = updateDistanceScore(player, x, y-i, axeTop, i, t)
		b = updateDistanceScore(player, x, y+i, axeBottom, i, b)
		// Axe lefttop rightbottom
		lt = updateDistanceScore(player, x-i, y-i, axeLeftTop, i, lt)
		rb = updateDistanceScore(player, x+i, y+i, axeRightBottom, i, rb)
		// Axe righttop leftbottom
		rt = updateDistanceScore(player, x+i, y-i, axeRightTop, i, rt)
		lb = updateDistanceScore(player, x-i, y+i, axeLeftBottom, i, lb)
	}
}
