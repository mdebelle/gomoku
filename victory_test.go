package main

import (
	"testing"
	"fmt"
)

var _ = fmt.Println

func cmpCaptures(a, b []Position) bool {
	if a == nil && b == nil {
		return true
	}

	aBoard := [19][19]bool{}
	bBoard := [19][19]bool{}

	fillBoard := func (board *[19][19]bool, capts []Position) {
		for _, pos := range capts {
			board[pos.y][pos.x] = true
		}
	}

	fillBoard(&aBoard, a)
	fillBoard(&bBoard, b)

	return aBoard == bBoard
}

func testVictory(align, alignTest AlignmentType, captures, capturesTest []Position, t *testing.T) {
	if align != alignTest {
		t.Error("\nExpected :", alignTest, "\nGot :", align)
	}
	if align == capturableAlignment && !cmpCaptures(captures, capturesTest) {
		t.Error("\nExpected :", capturesTest, "\nGot :", captures)
	}
}

// [1 2 3 4 .]
func TestVictory(t *testing.T) {
	board, free, align, capture := Board{}, [2]Board{}, [2]Board{}, [3]int{}
	checkRules(&board, &free, &align, &capture, 0, 0, player_one)
	checkRules(&board, &free, &align, &capture, 1, 0, player_one)
	checkRules(&board, &free, &align, &capture, 2, 0, player_one)
	checkRules(&board, &free, &align, &capture, 3, 0, player_one)
	alignType, captures := checkVictory(&board, 3, 0, player_one)
	testVictory(alignType, regularAlignment, captures, nil, t)
}

// [1 2 3 4 5]
func TestVictory2(t *testing.T) {
	board, free, align, capture := Board{}, [2]Board{}, [2]Board{}, [3]int{}
	checkRules(&board, &free, &align, &capture, 0, 0, player_one)
	checkRules(&board, &free, &align, &capture, 1, 0, player_one)
	checkRules(&board, &free, &align, &capture, 2, 0, player_one)
	checkRules(&board, &free, &align, &capture, 3, 0, player_one)
	checkRules(&board, &free, &align, &capture, 4, 0, player_one)
	alignType, captures := checkVictory(&board, 4, 0, player_one)
	testVictory(alignType, winningAlignment, captures, nil, t)
}

// [1 2 5 3 4]
func TestVictory3(t *testing.T) {
	board, free, align, capture := Board{}, [2]Board{}, [2]Board{}, [3]int{}
	checkRules(&board, &free, &align, &capture, 0, 0, player_one)
	checkRules(&board, &free, &align, &capture, 1, 0, player_one)
	checkRules(&board, &free, &align, &capture, 3, 0, player_one)
	checkRules(&board, &free, &align, &capture, 4, 0, player_one)
	checkRules(&board, &free, &align, &capture, 2, 0, player_one)
	alignType, captures := checkVictory(&board, 2, 0, player_one)
	testVictory(alignType, winningAlignment, captures, nil, t)
}

// [X . . . . .]
// [O O O O 0 .]
// [. . O . . .]
// [. . . . . .]
//        ^ There is a capture here !
func TestVictory4(t *testing.T) {
	board, free, align, capture := Board{}, [2]Board{}, [2]Board{}, [3]int{}
	checkRules(&board, &free, &align, &capture, 0, 1, player_one)
	checkRules(&board, &free, &align, &capture, 1, 1, player_one)
	checkRules(&board, &free, &align, &capture, 2, 1, player_one)
	checkRules(&board, &free, &align, &capture, 3, 1, player_one)
	checkRules(&board, &free, &align, &capture, 0, 0, player_two)
	checkRules(&board, &free, &align, &capture, 2, 2, player_one)
	checkRules(&board, &free, &align, &capture, 4, 1, player_one)
	alignType, captures := checkVictory(&board, 4, 1, player_one)
	testVictory(alignType, capturableAlignment, captures, []Position{{3, 3}}, t)
}

// [X X . X . .]
// [O O O O 0 .]
// [. O O . . .]
// [. . . . . .]
//  ^ ^   ^ Wow so many captures !
func TestVictory5(t *testing.T) {
	board, free, align, capture := Board{}, [2]Board{}, [2]Board{}, [3]int{}
	checkRules(&board, &free, &align, &capture, 0, 1, player_one)
	checkRules(&board, &free, &align, &capture, 1, 1, player_one)
	checkRules(&board, &free, &align, &capture, 2, 1, player_one)
	checkRules(&board, &free, &align, &capture, 3, 1, player_one)
	checkRules(&board, &free, &align, &capture, 2, 2, player_one)
	checkRules(&board, &free, &align, &capture, 1, 2, player_one)
	checkRules(&board, &free, &align, &capture, 0, 0, player_two)
	checkRules(&board, &free, &align, &capture, 1, 0, player_two)
	checkRules(&board, &free, &align, &capture, 3, 0, player_two)
	checkRules(&board, &free, &align, &capture, 4, 1, player_one)
	alignType, captures := checkVictory(&board, 4, 1, player_one)
	testCaptures := []Position{{0, 3}, {1, 3}, {3, 3}}
	testVictory(alignType, capturableAlignment, captures, testCaptures, t)
}

// [X . . . . .]
// [O O O 0 O O]
// [O . . . . .]
// [. . . . . .]
//  ^ You can capture me so what ? There is still five pawns aligned !
func TestVictory6(t *testing.T) {
	board, free, align, capture := Board{}, [2]Board{}, [2]Board{}, [3]int{}
	checkRules(&board, &free, &align, &capture, 0, 1, player_one)
	checkRules(&board, &free, &align, &capture, 0, 2, player_one)
	checkRules(&board, &free, &align, &capture, 1, 1, player_one)
	checkRules(&board, &free, &align, &capture, 2, 1, player_one)
	checkRules(&board, &free, &align, &capture, 4, 1, player_one)
	checkRules(&board, &free, &align, &capture, 5, 1, player_one)
	checkRules(&board, &free, &align, &capture, 0, 0, player_two)
	checkRules(&board, &free, &align, &capture, 3, 1, player_one)
	alignType, captures := checkVictory(&board, 3, 1, player_one)
	testCaptures := []Position{}
	testVictory(alignType, winningAlignment, captures, testCaptures, t)
}

// [X X . . . X X]
// [O O O 0 O O O]
// [O O . . . O O]
// [. . . . . . .]
//  ^ ^       ^ ^ Same here !
func TestVictory7(t *testing.T) {
	board, free, align, capture := Board{}, [2]Board{}, [2]Board{}, [3]int{}
	checkRules(&board, &free, &align, &capture, 0, 1, player_one)
	checkRules(&board, &free, &align, &capture, 0, 2, player_one)
	checkRules(&board, &free, &align, &capture, 1, 1, player_one)
	checkRules(&board, &free, &align, &capture, 1, 2, player_one)
	checkRules(&board, &free, &align, &capture, 2, 1, player_one)
	checkRules(&board, &free, &align, &capture, 4, 1, player_one)
	checkRules(&board, &free, &align, &capture, 5, 1, player_one)
	checkRules(&board, &free, &align, &capture, 5, 2, player_one)
	checkRules(&board, &free, &align, &capture, 6, 1, player_one)
	checkRules(&board, &free, &align, &capture, 6, 2, player_one)
	checkRules(&board, &free, &align, &capture, 0, 0, player_two)
	checkRules(&board, &free, &align, &capture, 1, 0, player_two)
	checkRules(&board, &free, &align, &capture, 5, 0, player_two)
	checkRules(&board, &free, &align, &capture, 6, 0, player_two)
	checkRules(&board, &free, &align, &capture, 3, 1, player_one)
	alignType, captures := checkVictory(&board, 3, 1, player_one)
	testCaptures := []Position{}
	testVictory(alignType, winningAlignment, captures, testCaptures, t)
}

// [X X . X . X X]
// [O O O 0 O O O]
// [O O . . . O O]
// [. . . . . . .]
//  ^           ^ Well not here though ...
func TestVictory8(t *testing.T) {
	board, free, align, capture := Board{}, [2]Board{}, [2]Board{}, [3]int{}
	checkRules(&board, &free, &align, &capture, 0, 1, player_one)
	checkRules(&board, &free, &align, &capture, 0, 2, player_one)
	checkRules(&board, &free, &align, &capture, 1, 1, player_one)
	checkRules(&board, &free, &align, &capture, 1, 2, player_one)
	checkRules(&board, &free, &align, &capture, 2, 1, player_one)
	checkRules(&board, &free, &align, &capture, 4, 1, player_one)
	checkRules(&board, &free, &align, &capture, 5, 1, player_one)
	checkRules(&board, &free, &align, &capture, 5, 2, player_one)
	checkRules(&board, &free, &align, &capture, 6, 1, player_one)
	checkRules(&board, &free, &align, &capture, 6, 2, player_one)
	checkRules(&board, &free, &align, &capture, 0, 0, player_two)
	checkRules(&board, &free, &align, &capture, 1, 0, player_two)
	checkRules(&board, &free, &align, &capture, 5, 0, player_two)
	checkRules(&board, &free, &align, &capture, 6, 0, player_two)
	checkRules(&board, &free, &align, &capture, 3, 0, player_two)
	checkRules(&board, &free, &align, &capture, 3, 1, player_one)
	alignType, captures := checkVictory(&board, 3, 1, player_one)
	testCaptures := []Position{{0, 3}, {6, 3}}
	testVictory(alignType, capturableAlignment, captures, testCaptures, t)
}

// [. . O . .]
// [O O 0 O O]
// [. . O . .]
// [. . O . .]
// [. . O . .]
//      ^ Christus Powaa
func TestVictory9(t *testing.T) {
	board, free, align, capture := Board{}, [2]Board{}, [2]Board{}, [3]int{}
	checkRules(&board, &free, &align, &capture, 0, 1, player_one)
	checkRules(&board, &free, &align, &capture, 1, 1, player_one)
	checkRules(&board, &free, &align, &capture, 3, 1, player_one)
	checkRules(&board, &free, &align, &capture, 4, 1, player_one)

	checkRules(&board, &free, &align, &capture, 2, 0, player_one)
	checkRules(&board, &free, &align, &capture, 2, 2, player_one)
	checkRules(&board, &free, &align, &capture, 3, 2, player_one)
	checkRules(&board, &free, &align, &capture, 4, 2, player_one)

	checkRules(&board, &free, &align, &capture, 2, 1, player_one)
	alignType, captures := checkVictory(&board, 2, 1, player_one)
	testCaptures := []Position{}
	testVictory(alignType, winningAlignment, captures, testCaptures, t)
}

// [X . O . .]
// [O O 0 O O]
// [. . O . .]
// [. . O .< Here you capture on both alignments
// [. . O . .]
func TestVictory10(t *testing.T) {
	board, free, align, capture := Board{}, [2]Board{}, [2]Board{}, [3]int{}
	checkRules(&board, &free, &align, &capture, 0, 1, player_one)
	checkRules(&board, &free, &align, &capture, 1, 1, player_one)
	checkRules(&board, &free, &align, &capture, 3, 1, player_one)
	checkRules(&board, &free, &align, &capture, 4, 1, player_one)

	checkRules(&board, &free, &align, &capture, 2, 0, player_one)
	checkRules(&board, &free, &align, &capture, 2, 2, player_one)
	checkRules(&board, &free, &align, &capture, 2, 3, player_one)
	checkRules(&board, &free, &align, &capture, 2, 4, player_one)

	checkRules(&board, &free, &align, &capture, 0, 0, player_two)

	checkRules(&board, &free, &align, &capture, 2, 1, player_one)

	alignType, captures := checkVictory(&board, 2, 1, player_one)
	testCaptures := []Position{{3, 3}}
	checkRules(&board, &free, &align, &capture, 1, 1, player_one)
	testVictory(alignType, capturableAlignment, captures, testCaptures, t)
}

// [. X O . .]
// [O O 0 O O]
// [. O O . .]
// [. . O . .]
// [. ^ O . .]
//   Does not work: Vertical alignment unbreakable
func TestVictory11(t *testing.T) {
	board, free, align, capture := Board{}, [2]Board{}, [2]Board{}, [3]int{}
	checkRules(&board, &free, &align, &capture, 0, 1, player_one)
	checkRules(&board, &free, &align, &capture, 1, 1, player_one)
	checkRules(&board, &free, &align, &capture, 1, 2, player_one)
	checkRules(&board, &free, &align, &capture, 3, 1, player_one)
	checkRules(&board, &free, &align, &capture, 4, 1, player_one)

	checkRules(&board, &free, &align, &capture, 2, 0, player_one)
	checkRules(&board, &free, &align, &capture, 2, 2, player_one)
	checkRules(&board, &free, &align, &capture, 2, 3, player_one)
	checkRules(&board, &free, &align, &capture, 2, 4, player_one)

	checkRules(&board, &free, &align, &capture, 1, 0, player_two)

	checkRules(&board, &free, &align, &capture, 2, 1, player_one)
	alignType, captures := checkVictory(&board, 2, 1, player_one)
	testCaptures := []Position{}
	testVictory(alignType, winningAlignment, captures, testCaptures, t)
}

// [. X O . .]
// [O O 0 O O]
// [. O O X .]
// [^ . O . .]
// [. ^ O . .]
//   Does not work either: Captures are on different positions
func TestVictory12(t *testing.T) {
	board, free, align, capture := Board{}, [2]Board{}, [2]Board{}, [3]int{}
	checkRules(&board, &free, &align, &capture, 0, 1, player_one)
	checkRules(&board, &free, &align, &capture, 1, 1, player_one)
	checkRules(&board, &free, &align, &capture, 1, 2, player_one)
	checkRules(&board, &free, &align, &capture, 3, 1, player_one)
	checkRules(&board, &free, &align, &capture, 4, 1, player_one)

	checkRules(&board, &free, &align, &capture, 2, 0, player_one)
	checkRules(&board, &free, &align, &capture, 2, 2, player_one)
	checkRules(&board, &free, &align, &capture, 2, 3, player_one)
	checkRules(&board, &free, &align, &capture, 2, 4, player_one)

	checkRules(&board, &free, &align, &capture, 1, 0, player_two)
	checkRules(&board, &free, &align, &capture, 3, 2, player_two)

	checkRules(&board, &free, &align, &capture, 2, 1, player_one)
	alignType, captures := checkVictory(&board, 2, 1, player_one)
	testCaptures := []Position{}
	testVictory(alignType, winningAlignment, captures, testCaptures, t)
}

// [. X O . .]
// [O O 0 O O]
// [. O O . .]
// [. . O O X]
// [. ^ O . .]
//   But this do. Different captures on different alignments
func TestVictory13(t *testing.T) {
	board, free, align, capture := Board{}, [2]Board{}, [2]Board{}, [3]int{}
	checkRules(&board, &free, &align, &capture, 0, 1, player_one)
	checkRules(&board, &free, &align, &capture, 1, 1, player_one)
	checkRules(&board, &free, &align, &capture, 1, 2, player_one)
	checkRules(&board, &free, &align, &capture, 3, 1, player_one)
	checkRules(&board, &free, &align, &capture, 4, 1, player_one)

	checkRules(&board, &free, &align, &capture, 2, 0, player_one)
	checkRules(&board, &free, &align, &capture, 2, 2, player_one)
	checkRules(&board, &free, &align, &capture, 2, 3, player_one)
	checkRules(&board, &free, &align, &capture, 3, 3, player_one)
	checkRules(&board, &free, &align, &capture, 2, 4, player_one)

	checkRules(&board, &free, &align, &capture, 1, 0, player_two)
	checkRules(&board, &free, &align, &capture, 4, 3, player_two)

	checkRules(&board, &free, &align, &capture, 2, 1, player_one)
	alignType, captures := checkVictory(&board, 2, 1, player_one)
	testCaptures := []Position{{1, 3}}
	testVictory(alignType, capturableAlignment, captures, testCaptures, t)
}
