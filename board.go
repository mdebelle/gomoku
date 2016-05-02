package main

type AIBoard struct {
	board		*Board
	freeThrees	*[2]Board
	capturesNb	*[3]int
	player		int
}

type Move struct {
	pos			Position
	captures	[]Position
	freeThrees	bool // Dunno. New free threes positions/axes.
}

func (board *AIBoard) SwitchPlayer() {
	board.player = -board.player
}

func (board *AIBoard) isValidMove(x, y int) bool {
	return isInBounds(x, y) &&
		board.board[y][x] == empty &&
		!doesDoubleFreeThree(board.freeThrees, x, y, board.player)
}

func (board *AIBoard) GetSearchSpace() []Position {
	moves := make([]Position, 0, 10)
	alreadyChecked := [19][19]bool {}

	checkAxis := func(x, y, incx, incy int) {
		if board.isValidMove(x + incx, y + incy) && !alreadyChecked[y + incy][x + incx] {
			alreadyChecked[y + incy][x + incx] = true
			moves = append(moves, Position{x + incx, y + incy})
		}
		if board.isValidMove(x - incx, y - incy) && !alreadyChecked[y - incy][x - incx] {
			alreadyChecked[y - incy][x - incx] = true
			moves = append(moves, Position{x - incx, y - incy})
		}
	}

	for y := 0; y < 19; y++ {
		for x := 0; x < 19; x++ {
			if board.board[y][x] != empty {
				for i := 0; i < 2; i++ {
					checkAxis(x, y, i, 0)
					checkAxis(x, y, 0, i)
					checkAxis(x, y, i, i)
					checkAxis(x, y, i, -i)
				}
			}
		}
	}

	return moves
}

func (board *AIBoard) DoMove(pos Position) []Position {
	captures := make([]Position, 0, 16)
	doMove(board.board, pos.x, pos.y, board.player, &captures)
	board.capturesNb[board.player + 1] += len(captures)
	return captures
}

func (board *AIBoard) UndoMove(pos Position, captures *[]Position) {
	undoMove(board.board, pos.x, pos.y, board.player, captures)
	board.capturesNb[board.player + 1] -= len(*captures)
}
