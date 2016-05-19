package main

import (
	"testing"
)

func test(free *[2]Board, sample *Board, color int, t *testing.T) {
	if (free[(color + 1) / 2] != *sample) {
		t.Error("Expected :\n", sample[0], "\nGot :\n", free[1][0])
	}
}

// [..O.O.]
func TestFreeThrees0(t *testing.T) {
	board, free, align, capture := Board{}, [2]Board{}, [2]Board{}, [3]int{}
	checkRules(&board, &free, &align, &capture, 1, 0, player_one)
	checkRules(&board, &free, &align, &capture, 3, 0, player_one)
	test(&free, &Board{{0, 0, 2, 0, 2}}, player_one, t)
}

// [.0.0.X]
func TestFreeThrees1(t *testing.T) {
	board, free, align, capture := Board{}, [2]Board{}, [2]Board{}, [3]int{}
	checkRules(&board, &free, &align, &capture, 3, 0, player_one)
	checkRules(&board, &free, &align, &capture, 1, 0, player_one)
	checkRules(&board, &free, &align, &capture, 5, 0, player_two)
	test(&free, &Board{}, player_one, t)
}

// [.0.0.X]
func TestFreeThrees2(t *testing.T) {
	board, free, align, capture := Board{}, [2]Board{}, [2]Board{}, [3]int{}
	checkRules(&board, &free, &align, &capture, 3, 0, player_one)
	checkRules(&board, &free, &align, &capture, 5, 0, player_two)
	checkRules(&board, &free, &align, &capture, 1, 0, player_one)
	test(&free, &Board{}, player_one, t)
}

// [.0.0.0..]
func TestFreeThrees3(t *testing.T) {
	board, free, align, capture := Board{}, [2]Board{}, [2]Board{}, [3]int{}
	checkRules(&board, &free, &align, &capture, 3, 0, player_one)
	checkRules(&board, &free, &align, &capture, 5, 0, player_one)
	checkRules(&board, &free, &align, &capture, 1, 0, player_one)
	test(&free, &Board{{0, 0, 0, 0, 2, 0, 2}}, player_one, t)
}

// [...00X]
func TestFreeThrees4(t *testing.T) {
	board, free, align, capture := Board{}, [2]Board{}, [2]Board{}, [3]int{}
	checkRules(&board, &free, &align, &capture, 3, 0, player_one)
	checkRules(&board, &free, &align, &capture, 4, 0, player_one)
	checkRules(&board, &free, &align, &capture, 5, 0, player_two)
	test(&free, &Board{}, player_one, t)
}

// Capture
// [..OXX0]
func TestFreeThrees5(t *testing.T) {
	board, free, align, capture := Board{}, [2]Board{}, [2]Board{}, [3]int{}
	checkRules(&board, &free, &align, &capture, 2, 0, player_one)
	checkRules(&board, &free, &align, &capture, 3, 0, player_two)
	checkRules(&board, &free, &align, &capture, 4, 0, player_two)
	checkRules(&board, &free, &align, &capture, 5, 0, player_one)
	test(&free, &Board{{0, 0, 0, 2, 2}}, player_one, t)
}
