// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   victories.go                                       :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: tmielcza <marvin@42.fr>                    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2016/05/18 19:28:48 by tmielcza          #+#    #+#             //
//   Updated: 2016/05/20 15:40:43 by tmielcza         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package main

func checkVictory(board *Board, capturesNb *[3]int, x, y, player int) (AlignmentType, []Position) {
	var captures []Position = nil

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

	checkAxis := func (incx, incy int) (AlignmentType, []Position) {
		countPawnsOnDir := func (incx, incy int) int {
			x, y := x + incx, y + incy
			i := 0
			for i < 4 && isInBounds(x, y) && board[y][x] == player {
				i++
				x += incx
				y += incy
			}
			return i
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

	getPossibleCounterCaptures := func () {
		for y := 0; y < 19; y++ {
			for x := 0; x < 19; x++ {
				if board[y][x] == player {
					testCapturability(x, y, &captures)
				}
			}
		}
	}

	updateCaptures(checkAxis(1, 0))
	updateCaptures(checkAxis(0, 1))
	updateCaptures(checkAxis(1, 1))
	updateCaptures(checkAxis(-1, 1))

	if best > regularAlignment && capturesNb[-player + 1] == 8 {
		getPossibleCounterCaptures()
		best = capturableAlignment
	}
	if best == capturableAlignment {
		if len(captures) == 0 {
			return winningAlignment, nil
		}
	}
	return best, captures
}
