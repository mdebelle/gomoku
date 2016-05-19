package main

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
