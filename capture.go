package main


// TODO: Create a list of possible anti-victory captures
func checkCaptures(values *Board, nb, x, y, incx, incy int) bool {
	checkAxis := func (x, y, incx, incy int) bool {
		if !isInBounds(x - incx, y - incy) || !isInBounds(x + 2 * incx, y + 2 * incy) {
			return false
		}
		if values[y + incy][x + incx] == nb {
			if values[y + 2 * incy][x + 2 * incx] == -nb && values[y - incy][x - incx] == 0 {
				victory.X = x - incx
				victory.Y = y - incy
				victory.Todo = true
				return true
			} else if values[y + 2 * incy][x + 2 * incx] == 0   && values[y - incy][x - incx] == -nb {
				victory.X = x + 2 * incx
				victory.Y = y + 2 * incy
				victory.Todo = true
				return true
			}
		}
		return false
	}
	f := func (incx, incy int) bool {
		x, y := x, y
		for i := 0; i < 5; i++ {
			if !isInBounds(x, y) || values[y][x] != nb {
				return false
			} else if checkAxis(x, y, -1, -1) || checkAxis(x, y, 1, 1) ||
				checkAxis(x, y, 1, -1) || checkAxis(x, y, -1, 1) ||
				checkAxis(x, y, 0, -1) || checkAxis(x, y, 0, 1) ||
				checkAxis(x, y, -1, 0) || checkAxis(x, y, 1, 0) {
				return true
			}
			x += incx
			y += incy
		}
		return false
	}
	return f(incx, incy) || f(-incx, -incy)
}

func getCaptures(board *Board, x, y, player  int, captures *[]Position) {
	captureOnAxis := func (incx, incy int) {
		if !isInBounds(x + 3 * incx, y + 3 * incy) {
			return
		}
		if	board[y + incy][x + incx] == -player &&
			board[y + 2 * incy][x + 2 * incx] == -player &&
		 	board[y + 3 * incy][x + 3 * incx] == player {
			*captures = append(*captures, Position{x + incx, y + incy})
			*captures = append(*captures, Position{x + incx * 2, y + incy * 2})
		}
	}
	captureOnAxis(-1, -1)
	captureOnAxis(1, 1)
	captureOnAxis(1, -1)
	captureOnAxis(-1, 1)
	captureOnAxis(0, -1)
	captureOnAxis(0, 1)
	captureOnAxis(-1, 0)
	captureOnAxis(1, 0)
}

func doCaptures(board *Board, captures *[]Position) {
	for _, capture := range *captures {
		board[capture.y][capture.x] = empty
	}
}

func undoCaptures(board *Board, captures *[]Position, player int) {
	for _, capture := range *captures {
		board[capture.y][capture.x] = -player
	}
}
