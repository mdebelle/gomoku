package main


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

func updateAlign(board *Board, alignTable *[2]Board, x, y, player int) []AlignScore {

	var lst []AlignScore
	
	clearOponentSituation := func (player, py, px, axe, start int) {
		for j := start; j < 5; j++ {
			alignTable[(-player+1)/2][py][px] &= ^(1 << uint(axe-j))
		}
	}

	updateDistanceScore := func (player, px, py, axe, i int, state bool) bool {
		if isInBounds(px, py) {
			if !state && board[py][px] == 0 {
				lst = append(lst, AlignScore{alignTable[(player_one+1)/2][py][px], alignTable[(player_two+1)/2][py][px], px, py})
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
	return lst
}
