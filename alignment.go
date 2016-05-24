// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   alignment.go                                       :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: tmielcza <marvin@42.fr>                    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2016/05/23 18:35:17 by tmielcza          #+#    #+#             //
//   Updated: 2016/05/24 15:47:53 by tmielcza         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package main

import (
	"fmt"
	"math"
)

var fmt_debug = fmt.Println

// return weight of alignement in the selected axis
// and time to expect an wining alignment
// if space missing returns +inf for time
func winingAlignement(board *Board, axe1, axe2, x, y, incx, incy, player int) (int, int){

	// AlignementGagnant
	if axe1 == 15 || axe2 == 15 || (axe1 == 14 && axe2 >= 8) || (axe2 == 14 && axe1 >= 8) || (axe1 >= 12 && axe2 >= 12) {
		return math.MaxInt32, 1
	}

	var t1, t2 [5]int
	spaceAndChaine := func (axe1, axe2, incx, incy int) {

		j1, j2 := 0, 0
		lock1, lock2 := false, false
		for i:= 0; i < 4; i++ {
			if !lock1 && isInBounds(x+(i*incx), y+(i*incy)) {
				if ((axe1 >> uint(i)) & 1) == 1 {
					if (j1 % 2 != 0) { j1++ }
					t1[j1]++
				} else  {
					if (j1 % 2 == 0) { j1++ }
					if board[y+(i*incy)][x+(i*incx)] == 0 {
						t1[j1]++
					} else {
						lock1 = true
					}
				}
			}
			if !lock2 && isInBounds(x-(i*incx), y-(i*incy)) {
				if ((axe2 >> uint(i)) & 1) == 1 {
					if (j2 % 2 != 0) { j2++ }
					t2[j2]++
				} else {
					if (j2 % 2 == 0) { j2++ }
					if board[y-(i*incy)][x-(i*incx)] == empty {
						t2[j2]++
					} else {
						lock2 = true
					}
				}
			}
		}
	}
	spaceAndChaine(axe1, axe2, incx, incy)

	chaine, space := 1, 1
	// Possibility of wining alignment
	for i := 0; i < 5; i++ {
		
		if i % 2 == 0 {
			chaine += (t1[i] + t2[i])
			if chaine + space >= 5 {
				return axe1+axe2, space
			}
		} else {
			if i < 4 {
				if t1[i] < t2[i] {
					if chaine + space + t1[i] + t1[i+1] >= 5 {
						return axe1+axe2, space + t1[i]
					}
				} else { 
					if chaine + space + t2[i] + t2[i+1] >= 5 {
						return axe1+axe2, space + t2[i]
					}
				}
				space += (t1[i]+t2[i])
			}
		}
	}

	// Space missing
	return axe1+axe2, math.MaxInt32
}

func getBestScore(board *Board, alignTable *[2]Board, x, y, player int) (int, int, int, int) {
	
	var s, t [4]int
	v := alignTable[(player+1)/2][y][x]

	applikmask := func(m, i int) int {
		return ((v & m) >> uint(i)) 
	}


	s[0], t[0] = winingAlignement(board, applikmask(0x0000000f, 0), applikmask(0x000000f0, 4), x, y, -1, 0, player)
	s[1], t[1] = winingAlignement(board, applikmask(0x00000f00, 8), applikmask(0x0000f000, 12), x, y, 0, -1, player)
	s[2], t[2] = winingAlignement(board, applikmask(0x000f0000, 16), applikmask(0x00f00000, 20), x, y, -1, -1, player)
	s[3], t[3] = winingAlignement(board, applikmask(0x0f000000, 24), applikmask(0xf0000000, 28), x, y, 1, -1, player)

	/*
	fmt.Printf("---\n")
	for i := 0; i < 4; i++ {
		fmt.Printf("score: %d time: %d\n", s[i], t[i])
	}
	//*/

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
	return s[indexmin], t[indexmin], s[indexmax], t[indexmax]
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

