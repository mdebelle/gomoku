// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   victories.go                                       :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: tmielcza <marvin@42.fr>                    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2016/05/18 19:28:48 by tmielcza          #+#    #+#             //
//   Updated: 2016/05/18 19:28:59 by tmielcza         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package main

func checkVictory2(board *Board, x, y, player int) (AlignmentType, []Position) {
	var captures []Position = nil

	checkAxis := func (incx, incy int) (AlignmentType, []Position) {
		countPawnsOnDir := func (incx, incy int) int {
			x, y := x + incx, y + incy
			i := 0
			for ; i < 4 && isInBounds(x, y) && board[y][x] == player; {
				i++
				x += incx
				y += incy
			}
			return i
		}

		testCapturability := func (x, y int, captures *[]Position) bool {
			checkCaptureAxis := func (incx, incy int) bool {
				if !isInBounds(x - incx, y - incy) || !isInBounds(x + 2 * incx, y + 2 * incy) {
					return false
				} else if board[y + incy][x + incx] == player {
					if board[y + 2 * incy][x + 2 * incx] == -player && board[y - incy][x - incx] == 0 {
						*captures = append(*captures, Position{x - incx, y - incy})
						return true
					} else if board[y + 2 * incy][x + 2 * incx] == 0 && board[y - incy][x - incx] == -player {
						*captures = append(*captures, Position{x + incx * 2, y + incy * 2})
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

		getBreakingCaptures := func (x, y, incx, incy, pawns int) (bool, []Position) {
			captures := make([]Position, 0, 1)

			capturable := false
			start := pawns - 4
			end := 4
			x, y = x + incx * start, y + incy * start
			for i := start; i <= end; i++ {
				capturable = testCapturability(x, y, &captures) || capturable
				x += incx
				y += incy
			}
			return capturable, captures
		}

		pawnsRight := countPawnsOnDir(incx, incy)
		pawnsLeft := countPawnsOnDir(-incx, -incy)
		pawns := pawnsRight + pawnsLeft
		if pawns >= 4 {
			isBreakable, captures := getBreakingCaptures(x - incx * pawnsLeft, y - incy * pawnsLeft, incx, incy, pawns)
			if isBreakable {
				return capturableAlignment, captures
			}
			return winningAlignment, nil
		}
		return regularAlignment, nil
	}

	best := regularAlignment

	updateCaptures := func (alignType AlignmentType, capts []Position) {
		if alignType > best {
			best = alignType
		}
		if alignType == capturableAlignment {
			if captures == nil {
				captures = capts
			} else {
				updatedCaptures := make([]Position, 0, len(captures))
				for _, pos := range(captures) {
					for _, otherPos := range(capts) {
						if pos == otherPos {
							updatedCaptures = append(updatedCaptures, pos)
							break
						}
					}
				}
				captures = updatedCaptures
			}
		}
	}

	updateCaptures(checkAxis(1, 0))
	updateCaptures(checkAxis(0, 1))
	updateCaptures(checkAxis(1, 1))
	updateCaptures(checkAxis(-1, 1))
	if best == capturableAlignment && len(captures) == 0 {
		return winningAlignment, nil
	}
	return best, captures
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
