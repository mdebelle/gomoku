package main

import (
	"time"
)

func checkDoubleThree(board, freeThrees *Board, x, y, color int) {
	// TODO: SOOOOO SLOOOOOW do something

	if debugFlag { defer timeFunc(time.Now(), "checkDoubleThree") }

	const (
		pat1 = 0x1A5 // -00--
		pat2 = 0x199 // -0-0-
		pat3 = 0x169 // --00-
		mask = 0x3FF // 2 * 5 bits
	)

	checkAxis := func(x, y, incx, incy, axis int) {
		if !isInBounds(x, y) || board[y][x] != empty {
			return
		}
		flags := uint32(0)

		// TODO: Perform this formatting on the entire line/column/diagonal
		tmp_x, tmp_y := x - incx*4, y - incy*4
		i := uint(0)
		for ; i < 4; i++ {
			if isInBounds(tmp_x, tmp_y) {
				flags |= uint32(board[tmp_y][tmp_x] * color + 1) << ((7 - i)*2)
			}
			tmp_x, tmp_y = tmp_x + incx, tmp_y + incy
		}
		tmp_x, tmp_y = x + incx, y + incy
		for ; i < 8; i++ {
			if isInBounds(tmp_x, tmp_y) {
				flags |= uint32(board[tmp_y][tmp_x] * color + 1) << ((7 - i)*2)
			} else {
				break
			}
			tmp_x, tmp_y = tmp_x + incx, tmp_y + incy
		}

		pos1 := (flags >> (2*3)) & mask
		pos2 := (flags >> (2*2)) & mask
		pos3 := (flags >> (2*1)) & mask
		pos4 := flags & mask
		createsFreeThree := pos4 == pat1 ||
							pos4 == pat2 ||
							pos4 == pat3 ||
							pos3 == pat1 ||
							pos3 == pat2 ||
							pos3 == pat3 ||
							pos2 == pat1 ||
							pos2 == pat2 ||
							pos2 == pat3 ||
							pos1 == pat1 ||
							pos1 == pat2 ||
							pos1 == pat3
		if createsFreeThree {
			freeThrees[y][x] |= axis
		} else {
			freeThrees[y][x] &= ^axis
		}
	}

	if board[y][x] == empty {
		checkAxis(x, y, 0, 1, VerticalAxisMask)
		checkAxis(x, y, 1, 0, HorizontalAxisMask)
		checkAxis(x, y, 1, -1, LeftDiagAxisMask)
		checkAxis(x, y, 1, 1, RightDiagAxisMask)
	}
	for i := 1; i <= 4; i++ {
		checkAxis(x, y + i, 0, 1, VerticalAxisMask)
		checkAxis(x, y - i, 0, 1, VerticalAxisMask)
		checkAxis(x + i, y, 1, 0, HorizontalAxisMask)
		checkAxis(x - i, y, 1, 0, HorizontalAxisMask)
		checkAxis(x + i, y - i, 1, -1, LeftDiagAxisMask)
		checkAxis(x - i, y + i, 1, -1, LeftDiagAxisMask)
		checkAxis(x + i, y + i, 1, 1, RightDiagAxisMask)
		checkAxis(x - i, y - i, 1, 1, RightDiagAxisMask)
	}
}

func doesDoubleFreeThreePlayer(freeThrees *Board, x, y int) bool {
	freeThreesCount := 0
	point := freeThrees[y][x]
	if (point == 0) {
		return false
	}
	for i := uint(0); i < 4; i++ {
		if (point & (1 << i)) != 0 {
			freeThreesCount++
		}
	}
	return freeThreesCount == 2
}

func doesDoubleFreeThree(freeThrees *[2]Board, x, y, player int) bool {
	playerId := (player + 1) / 2
	return doesDoubleFreeThreePlayer(&freeThrees[playerId], x, y)
}

func updateFreeThrees(board *Board, freeThrees *[2]Board, x, y, player int, captures []Position) {
	// TODO: Two functions -> update from move and update from move cancelation
	if debugFlag { defer timeFunc(time.Now(), "updateFreeThree") }

	if (board[y][x] != empty) {
		freeThrees[0][y][x] = 0
		freeThrees[1][y][x] = 0
	}
	checkDoubleThree(board, &freeThrees[(player + 1) / 2], x, y, player)
	checkDoubleThree(board, &freeThrees[(-player + 1) / 2], x, y, -player)
	for _, pos := range captures {
		if (board[pos.y][pos.x] != empty) {
			freeThrees[0][pos.y][pos.x] = 0
			freeThrees[1][pos.y][pos.x] = 0
		}
		checkDoubleThree(board, &freeThrees[(player + 1) / 2], pos.x, pos.y, player)
		checkDoubleThree(board, &freeThrees[(-player + 1) / 2], pos.x, pos.y, -player)
	}
}
