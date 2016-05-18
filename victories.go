package main

func checkVictory2(board *Board, x, y, player int) (AlignmentType, []Position) {
	var captures []Position = nil

	checkAxis := func (incx, incy int) AlignmentType {
		countPawnsOnDir := func (incx, incy int) int {
			x, y := x + incx, y + incy
			for i := 0; i < 4; i++ {
				if !isInBounds(x, y) || board[y][x] != player {
					return i
				}
				x += incx
				y += incy
			}
			return 4
		}

		testCapturability := func (x, y int) bool {
			checkCaptureAxis := func (incx, incy int) bool {
				if !isInBounds(x - incx, y - incy) || !isInBounds(x + 2 * incx, y + 2 * incy) {
					return false
				}
				if board[y + incy][x + incx] == player {
					if board[y + 2 * incy][x + 2 * incx] == -player && board[y - incy][x - incx] == 0 {
						captures = append(captures, Position{x - incx, y - incy})
						return true
					} else if board[y + 2 * incy][x + 2 * incx] == 0 && board[y - incy][x - incx] == -player {
						captures = append(captures, Position{x + incx * 2, y + incy * 2})
						return true
					}
				}
				return false
			}

			horizontal := checkCaptureAxis(1, 0) || checkCaptureAxis(-1, 0)
			vertical := checkCaptureAxis(0, 1) || checkCaptureAxis(0, -1)
			diagLeft := checkCaptureAxis(1, 1) || checkCaptureAxis(-1, -1)
			diagRight := checkCaptureAxis(-1, 1) || checkCaptureAxis(1, -1)
			return horizontal || vertical || diagLeft || diagRight
		}

		getBreakingCaptures := func (x, y, incx, incy, pawns int) bool {
			capturable := false
			start := pawns - 4
			end := 4
			x, y = x + incx * start, y + incy * start
			for i := start; i <= end; i++ {
				capturable = testCapturability(x, y) || capturable
				x += incx
				y += incy
			}
			return capturable
		}

		pawnsRight := countPawnsOnDir(incx, incy)
		pawnsLeft := countPawnsOnDir(-incx, -incy)
		pawns := pawnsRight + pawnsLeft
		if pawns >= 4 {
			if captures == nil {
				captures = make([]Position, 0, 1)
			}
			isBreakable := getBreakingCaptures(x - incx * pawnsLeft, y - incy * pawnsLeft, incx, incy, pawns)
			if !isBreakable {
				return winningAlignment
			} else {
				return capturableAlignment
			}
		} else {
			return regularAlignment
		}
	}

	getMax := func (axes []AlignmentType) AlignmentType {
		max := regularAlignment
		for _, axis := range(axes) {
			if axis > max {
				max = axis
			}
		}
		return max
	}

	victoryHorizontal := checkAxis(1, 0)
	victoryVertical := checkAxis(0, 1)
	victoryDiagLeft := checkAxis(1, 1)
	victoryDiagRight := checkAxis(-1, 1)
	max := getMax([]AlignmentType{victoryHorizontal, victoryVertical, victoryDiagLeft, victoryDiagRight})
	return max, captures
}

func checkVictory(values *Board, nb int, y int, x int) bool {
	f := func (incx, incy int) int {
		x, y := x + incx, y + incy
		for i := 0; i < 4; i++ {
			if !isInBounds(x, y) || values[y][x] != nb {
				return i
			}
			x += incx
			y += incy
		}
		return 5
	}
	if 	(f(-1, -1) + f(1, 1) >= 4 && !checkCaptures(values, nb, x, y, 1, 1)) ||
		(f(1, -1) + f(-1, 1) >= 4 && !checkCaptures(values, nb, x, y, 1, -1)) ||
		(f(0, -1) + f(0, 1) >= 4 && !checkCaptures(values, nb, x, y, 0, 1)) ||
		(f(-1, 0) + f(1, 0) >= 4 && !checkCaptures(values, nb, x, y, 1, 0)) {
		return true
	}
	return false
}
