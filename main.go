package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"time"
	"log"
)

const (
	empty = 0
	player_one = 1
	player_two = -player_one
	searchMaxTime = 500000000 * time.Nanosecond
	searchMaxDepth = 20 
)

const (
	VerticalAxis = 1 << iota
	HorizontalAxis
	LeftDiagAxis	// haut droite
	RightDiagAxis	// bas droite
)

type mustdo struct {
	Todo	bool
	X, Y    int
}

type Board [19][19]int
type FreeThreesAxis [2][19][19]int

var (
	winTitle string = "Go-Gomoku"
	winWidth, winHeight int = 800, 880
)

var victoir mustdo

func checkBounds(x, y int) bool {
	return x >= 0 && y >= 0 && x < 19 && y < 19
}

func checkVictory(values *Board, nb int, y int, x int) bool {
	f := func (incx, incy int) int {
		x, y := x + incx, y + incy
		for i := 0; i < 4; i++ {
			if !checkBounds(x, y) || values[y][x] != nb {
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

func checkCaptures(values *Board, nb, x, y, incx, incy int) bool {
	checkAxis := func (x, y, incx, incy int) bool {
		if !checkBounds(x - incx, y - incy) || !checkBounds(x + 2 * incx, y + 2 * incy) {
			return false
		}
		if values[y + incy][x + incx] == nb {
			if values[y + 2 * incy][x + 2 * incx] == -nb && values[y - incy][x - incx] == 0 {
				victoir.X = x - incx
				victoir.Y = y - incy
				victoir.Todo = true
				return true
			} else if values[y + 2 * incy][x + 2 * incx] == 0   && values[y - incy][x - incx] == -nb {
				victoir.X = x + 2 * incx
				victoir.Y = y + 2 * incy
				victoir.Todo = true
				return true
			}
		}
		return false
	}
	f := func (incx, incy int) bool {
		x, y := x + incx, y + incy
		for i := 0; i < 4; i++ {
			if !checkBounds(x, y) || values[y][x] != nb {
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

func checkDoubleThree(values, freeThrees *Board, x, y, color int) {

	checkAxis2 := func(x, y, incx, incy, axis int) {

		if !checkBounds(x, y) { return }
		
		if values[y][x] != 0 { 
			freeThrees[y][x] = 0
			return
		}

		if !checkBounds(x - incx, y - incy) || !checkBounds(x + incx, y + incy) { return }
		if values[y - incy][x - incx] == -color || values[y + incy][x + incx] == -color { 
			freeThrees[y][x] &= ^axis
			return
		}

		if checkBounds(x - (4 * incx), y - (4 * incy)) && 
			values[y - (4 * incy)][x - (4 * incx)] == 0 && values[y + (1 * incy)][x + (1 * incx)] == 0 {
			if values[y - (3 * incy)][x - (3 * incx)] == color {
				if values[y - (2 * incy)][x - (2 * incx)] == color &&
		   			values[y - (1 * incy)][x - (1 * incx)] == 0 {
		   				freeThrees[y][x] |= axis
		   				return
				}
				if values[y - (2 * incy)][x - (2 * incx)] == 0 &&
		   			values[y - (1 * incy)][x - (1 * incx)] == color {
		   				freeThrees[y][x] |= axis
		   				return
				}
			}
			if values[y - (3 * incy)][x - (3 * incx)] == 0 &&
		   		values[y - (2 * incy)][x - (2 * incx)] == color &&
		   		values[y - (1 * incy)][x - (1 * incx)] == color {
		   			freeThrees[y][x] |= axis
		   			return
			}
		}
		
		if checkBounds(x - (3 * incx), y - (3 * incy)) && checkBounds(x + (2 * incx), y + (2 * incy)) &&
			values[y - (3 * incy)][x - (3 * incx)] == 0 && values[y + (2 * incy)][x + (2 * incx)] == 0 {
			if values[y - (2 * incy)][x - (2 * incx)] == color {
				if values[y - (1 * incy)][x - (1 * incx)] == color &&
		   			values[y + (1 * incy)][x + (1 * incx)] == 0 {
		   				freeThrees[y][x] |= axis
		   				return
				}
				if values[y - (1 * incy)][x - (1 * incx)] == 0 &&
		   			values[y + (1 * incy)][x + (1 * incx)] == color {
		   				freeThrees[y][x] |= axis
		   				return
				}
			}
			if values[y - (2 * incy)][x - (2 * incx)] == 0 &&
		   		values[y - (1 * incy)][x - (1 * incx)] == color &&
		   		values[y + (1 * incy)][x + (1 * incx)] == color {
		   			freeThrees[y][x] |= axis
		   			return
			}
		}
		
		if checkBounds(x - (2 * incx), y - (2 * incy)) && checkBounds(x + (3 * incx), y + (3 * incy)) &&
			values[y - (2 * incy)][x - (2 * incx)] == 0 && values[y + (3 * incy)][x + (3 * incx)] == 0 {
			if values[y - (1 * incy)][x - (1 * incx)] == color {
				if values[y + (1 * incy)][x + (1 * incx)] == color &&
		   			values[y + (2 * incy)][x + (2 * incx)] == 0 {
		   				freeThrees[y][x] |= axis
		   				return
				}
				if values[y + (1 * incy)][x + (1 * incx)] == 0 &&
		   			values[y + (2 * incy)][x + (2 * incx)] == color {
		   				freeThrees[y][x] |= axis
		   				return
				}
			}
			if values[y - (1 * incy)][x - (1 * incx)] == 0 &&
		   		values[y + (1 * incy)][x + (1 * incx)] == color &&
		   		values[y + (2 * incy)][x + (2 * incx)] == color {
		   			freeThrees[y][x] |= axis
		   			return
			}
		}
		
		if checkBounds(x + (4 * incx), y + (4 * incy)) &&
			values[y - (1 * incy)][x - (1 * incx)] == 0 && values[y + (4 * incy)][x + (4 * incx)] == 0 {
			if values[y + (1 * incy)][x + (1 * incx)] == color {
				if values[y + (2 * incy)][x + (2 * incx)] == color &&
		   			values[y + (3 * incy)][x + (3 * incx)] == 0 {
		   				freeThrees[y][x] |= axis
		   				return
				}
				if values[y + (2 * incy)][x + (2 * incx)] == 0 &&
		   			values[y + (3 * incy)][x + (3 * incx)] == color {
		   				freeThrees[y][x] |= axis
		   				return
				}
			}
			if values[y + (1 * incy)][x + (1 * incx)] == 0 &&
		   		values[y + (2 * incy)][x + (2 * incx)] == color &&
		   		values[y + (3 * incy)][x + (3 * incx)] == color {
		   			freeThrees[y][x] |= axis
		   			return
			}
		}
		freeThrees[y][x] &= ^axis
		return
	}

	const (
		p_mine = iota
		p_theirs
		p_empty
		p_checked
	)

	const (
		s_start = iota
		s_1
		s_2
		s_3
		s_4
		s_5
		s_6
		s_7
		s_8
		s_9
		s_10
		s_11
		s_12
		s_13
		s_14
		s_15
		s_16
		s_end
		s_error
	)

	stateTable := [...][4]int {
//		    •    |   O    |	  Ø    |    @
		{ s_start, s_start, s_3,     s_error }, // start
		{ s_16   , s_error, s_error, s_8     }, // 1
		{ s_error, s_error, s_1,     s_16    }, // 2
		{ s_11,    s_start, s_4,     s_14    }, // 3
		{ s_9,     s_start, s_4,     s_5     }, // 4
		{ s_10,    s_error, s_6,     s_error }, // 5
		{ s_7,     s_error, s_error, s_error }, // 6
		{ s_8,     s_error, s_error, s_error }, // 7
		{ s_error, s_error, s_end,   s_error }, // 8
		{ s_13,    s_error, s_12,    s_10    }, // 9
		{ s_8,     s_error, s_7,     s_error }, // 10
		{ s_2,     s_start, s_12,    s_15    }, // 11
		{ s_1,     s_start, s_error, s_7     }, // 12
		{ s_error, s_start, s_1,     s_8     }, // 13
		{ s_15,    s_error, s_6,     s_error }, // 14
		{ s_16,    s_error, s_7,     s_error }, // 15
		{ s_error, s_error, s_8,     s_error }, // 16
		{ s_end,   s_end,   s_end,   s_end   }, // end
		{ s_error, s_error, s_error, s_error },
	}

	checkAxis3 := func(x, y, incx, incy int, axis int) {

		if !checkBounds(x, y) || values[y][x] != empty {
			return
		}
	
		state := s_start
		tmp_x, tmp_y := x - incx*4, y - incy*4
		for i := 0; i < 9; i++ {
			input := 0
			if !checkBounds(tmp_x, tmp_y) {
				input = p_theirs
			} else if tmp_y == y && tmp_x == x {
				input = p_checked
			} else {
				pos := values[tmp_y][tmp_x]
				if pos == color {
					input = p_mine
				} else if pos == -color {
					input = p_theirs
				} else {
					input = p_empty
				}
			}
			/*
			if state == s_error {
				freeThrees[y][x] &= ^axis
				return
			}
			*/
			state = stateTable[state][input]
			tmp_x += incx
			tmp_y += incy
		}
		if state == s_end {
			freeThrees[y][x] |= axis
		} else {
			freeThrees[y][x] &= ^axis
		}

	}

	const (
		pat1 = 0x1A5 // -00--
		pat2 = 0x199 // -0-0-
		pat3 = 0x169 // --00-
		mask = 0x3FF
	)

	checkAxis := func(x, y, incx, incy, axis int) {
		if !checkBounds(x, y) || values[y][x] != empty {
			return
		}
		flags := uint32(0)
		tmp_x, tmp_y := x - incx*4, y - incy*4
		for i := uint(0); i < 8; i++ {
			if !checkBounds(tmp_x, tmp_y) {
			} else if tmp_x == x && tmp_y == y {
				tmp_x += incx
				tmp_y += incy
				i--
				continue
			} else {
				flags |= uint32(values[tmp_y][tmp_x] * color + 1) << ((7 - i)*2)
			}
			tmp_x += incx
			tmp_y += incy
		}
		if  ((flags >> (2*3)) & mask) == pat1 ||
			((flags >> (2*3)) & mask) == pat2 ||
			((flags >> (2*3)) & mask) == pat3 ||
			((flags >> (2*2)) & mask) == pat1 ||
			((flags >> (2*2)) & mask) == pat2 ||
			((flags >> (2*2)) & mask) == pat3 ||
			((flags >> (2*1)) & mask) == pat1 ||
			((flags >> (2*1)) & mask) == pat2 ||
			((flags >> (2*1)) & mask) == pat3 ||
			((flags >> (2*0)) & mask) == pat1 ||
			((flags >> (2*0)) & mask) == pat2 ||
			((flags >> (2*0)) & mask) == pat3 {
			freeThrees[y][x] |= axis
		} else {
			freeThrees[y][x] &= ^axis
		}
	}

	//i ->  4 3 2 1   0   1 2 3 4
	//     | | | | | x,y | | | | |

	//     { |o|o| |  -  | } | | |
	//     { |o| |o|  -  | } | | |
	//     { | |o|o|  -  | } | | |

	//     | { |o|o|  -  | | } | |	
	//     | { |o| |  -  |o| } | |
	//     | { | |o|  -  |o| } | |

	//     | | { |o|  -  |o| | } |
	//     | | { |o|  -  | |o| } |
	//     | | { | |  -  |o|o| } |

	//     | | | { |  -  |o|o| | }
	//     | | | { |  -  |o| |o| }
	//     | | | { |  -  | |o|o| }

	if false {checkAxis2(0, 0, 0, 0, 0)}
	if false {checkAxis3(0, 0, 0, 0, 0)}
	if false {checkAxis(0, 0, 0, 0, 0)}

	for i := 1; i <= 4; i++ {
		checkAxis(x, y + i, 0, 1, VerticalAxis)
		checkAxis(x, y - i, 0, 1, VerticalAxis)
		checkAxis(x + i, y, 1, 0, HorizontalAxis)
		checkAxis(x - i, y, 1, 0, HorizontalAxis)
		checkAxis(x + i, y - i, 1, -1, LeftDiagAxis)
		checkAxis(x - i, y + i, 1, -1, LeftDiagAxis)
		checkAxis(x + i, y + i, 1, 1, RightDiagAxis)
		checkAxis(x - i, y - i, 1, 1, RightDiagAxis)
	}
		
	return
}

func doCaptures(values *Board, nb int, y int, x int) int {
	forcapture := func (incx, incy int) int {
		if !checkBounds(x + 3 * incx, y + 3 * incy) {
			return 0
		}
		if	values[y + incy][x + incx] == -nb &&
			values[y + 2 * incy][x + 2 * incx] == -nb &&
		 	values[y + 3 * incy][x + 3 * incx] == nb {
				values[y + incy][x + incx] = 0
				values[y + 2 * incy][x + 2 * incx] = 0
				return 2
		}
		return 0
	}
	return  forcapture(-1, -1) + forcapture(1, 1) +
			forcapture(1, -1) + forcapture(-1, 1) +
			forcapture(0, -1) + forcapture(0, 1) +
			forcapture(-1, 0) + forcapture(1, 0)
}

func checkRules(values *Board, freeThrees *[2]Board, capture *[3]int, x, y, player int) int {
	freeThreesCount := 0
	playerId := (player + 1) / 2
	for i := uint(0); i < 4; i++ {
		if (freeThrees[playerId][y][x] & (1 << i)) != 0 {
			freeThreesCount++
		}
	}
	if freeThreesCount == 2 {
		fmt.Printf("Nope\n")
		return player
	}
	values[y][x] = player
	freeThrees[0][y][x] = 0
	freeThrees[1][y][x] = 0
	checkDoubleThree(values, &freeThrees[playerId], x, y, player)
	checkDoubleThree(values, &freeThrees[(-player + 1) / 2], x, y, -player)
	victory := checkVictory(values, player, y, x)
	if victory == true {
		fmt.Printf("Victoire \\o/ %d\n", player)
		return 0
	}
	capture[player + 1] += doCaptures(values, player, y, x)
	if capture[player + 1] >= 10 {
		fmt.Printf("capture de ouf \\o/ %d\n", player)
		return 0
	}
	return -player
}

func mousePositionToGrid(val float64) int {
	t := int(math.Floor((val - 20) / 40))
	if t < 0 {
		t = 0
	} else if t > 18{
		t = 18
	}
	return t
}

func gridAnalyse(values *Board, nb int) (int, int) {
	f := func (incx, incy , x, y, nb int) int {
		x, y = x + incx, y + incy
		for i := 0; i < 4; i++ {
			if !checkBounds(x, y) || values[y][x] != nb {
				return i
			}
			x += incx
			y += incy
		}
		return 5
	}
	betterx, bettery := -1, -1
	max := 0
	for i := 0; i < 19; i++ {
		for j := 0; j < 19; j++ {
			if (values[i][j] == 0) {
				var t int
				tmp := f(-1, -1, j, i, player_two) + f(1, 1, j, i, player_two)
				t = f(1, -1, j, i, player_two) + f(-1, 1, j, i, player_two)
				if t > tmp {
					tmp = t
				}
				t = f(0, -1, j, i, player_two) + f(0, 1, j, i, player_two)
				if t > tmp {
					tmp = t
				}
				t = f(-1, 0, j, i, player_two) + f(1, 0, j, i, player_two)
				if t > tmp {
					tmp = t
				}
				t = (f(-1, -1, j, i, player_one) + f(1, 1, j, i, player_one)) * 2
				if t > tmp {
					tmp = t
				}
				t = (f(1, -1, j, i, player_one) + f(-1, 1, j, i, player_one)) * 2
				if t > tmp {
					tmp = t
				}
				t = (f(0, -1, j, i, player_one) + f(0, 1, j, i, player_one)) * 2
				if t > tmp {
					tmp = t
				}
				t = (f(-1, 0, j, i, player_one) + f(1, 0, j, i, player_one)) * 2
				if t > tmp {
					tmp = t
				}
				if tmp > max {
					max = tmp
					bettery = i
					betterx = j
				}
			}
		}
	}
	if (betterx < 0) {
		betterx = rand.Int() % 19
		bettery = rand.Int() % 19
	}
	return betterx, bettery
}

func init() {
    runtime.LockOSThread()
}

func evaluateAllBoard(player int, value *Board, better *[19][19][5]int, capture *[3]int) {
	for y := 0; y < 19; y++ {
		for x := 0; x < 19; x++ {
			evaluateBoard(value, x, y, player, better, capture)
		}
	}
}

func run() int {
	var event sdl.Event
	var running bool
	var err error
	var player int
	victoir.Todo = false
	var capture [3]int
	var values Board
	var freeThrees [2]Board
	var better [19][19][5]int

	f, err := os.OpenFile("testlogfile", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create file: %s\n", err)
		return 1
	}
	defer f.Close()

	log.SetOutput(f)
	log.Printf("---NEW GAME---\n")



	sdl.Init(sdl.INIT_EVERYTHING)

	window, err := sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 1
	}
	defer window.Destroy()
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return 2
	}
	defer renderer.Destroy()

	windowb, err := sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 1
	}
	defer windowb.Destroy()
	rendererb, err := sdl.CreateRenderer(windowb, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return 2
	}
	defer rendererb.Destroy()

	_= rendererb.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	running = true
	player = 1
	var px, py int
	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseButtonEvent:
				fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n", t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
				if  player == player_one &&  t.Type == 1025 {
					py = mousePositionToGrid(float64(t.Y))
					px = mousePositionToGrid(float64(t.X))
					fmt.Printf("Player -> x[%d] y [%d]\n", px, py)
					log.Printf("p1 -> X |%3d| Y|%3d|\n", px, py)
					if victoir.Todo == true {
						if px == victoir.X && py == victoir.Y {
							player = checkRules(&values, &freeThrees, &capture, px, py, player)
							victoir.Todo = false
						} else {
							fmt.Printf("you must play in [%d][%d]\n", victoir.X, victoir.Y)
						}
					} else if values[py][px] == 0 {
						player = checkRules(&values, &freeThrees, &capture, px, py, player)
					}
					evaluateAllBoard(player, &values, &better, &capture)
				}
			case *sdl.KeyUpEvent:
				fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n", t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
			}
		}
//*
		if player == player_two {
			if victoir.Todo == true {
				fmt.Printf("IA must play -> x[%d] y [%d]\n", victoir.X, victoir.Y)
				log.Printf("IA -> X |%3d| Y|%3d|\n", victoir.X, victoir.Y)
				player = checkRules(&values, &freeThrees, &capture, victoir.X, victoir.Y, player)
				victoir.Todo = false	
			} else {
				var x, y int
				x, y, better = search(&values, &freeThrees, player, px, py, 4, &capture)
				fmt.Printf("IA -> x[%d] y [%d]\n", x, y)
				log.Printf("IA -> X |%3d| Y|%3d|\n", x, y)
				if values[y][x] == 0 {
					player = checkRules(&values, &freeThrees, &capture, x, y, player)
				}
			}
		}
//*/
		_ = renderer.SetDrawColor(236, 240, 241, 0)
		renderer.Clear()
		drawGrid(renderer)
		drawClic(renderer, &values, &capture, &freeThrees[(-player + 1) / 2])
		renderer.Present()

		_ = rendererb.SetDrawColor(236, 240, 241, 0)
		rendererb.Clear()
		drawGrid(rendererb)
		draweval(rendererb, &better)
		rendererb.Present()
	}
	return 0
}

func main() {
	fmt.Printf("%d %d", HorizontalAxis, LeftDiagAxis)
	os.Exit(run())
}
