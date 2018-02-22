package pazudoraer

import (
	"testing"
)

func TestUpdateOnceNullBoard(t *testing.T) {
	b, _ := NewBoard(5, 6, "000000000000000000000000000000")
	config := &Config{3, 10}
	imp, _ := NewBoardImp(b, config)
	if imp.UpdateOnce() {
		t.Errorf("%v", b)
	}
	if imp.statics.combo != 0 {
		t.Errorf("%v", b)
	}
}

func TestUpdateOnce(t *testing.T) {
	b, _ := NewBoard(5, 6, "000000000000000000000000000000")
	config := &Config{3, 10}
	cases := []struct {
		board        string
		updatedBoard string
		combo        int
	}{
		{"111000000000000000000000000000", "000000000000000000000000000000", 1},
		{"100000100000100000000000000000", "000000000000000000000000000000", 1},
		{"000000000000000001000001000111", "000000000000000000000000000000", 1},
		{"000000000010010010010010011110", "000000000000000000000000000000", 1},
		{"000000000000000010000111000011", "000000000000000000000000000001", 1},
		{"333444211115212215211115666666", "000000000000000000000000002200", 6},
		{"566363644363422363233563111355", "000000566000644000422500233355", 4},
		{"245451243441233341243441222111", "000000000000040400040400045450", 4},
		{"111111111111111111111111111111", "000000000000000000000000000000", 1},
		{"111222222111111222222111111222", "000000000000000000000000000000", 10},
	}
	for _, c := range cases {
		b.SetBoard(c.board)
		imp, _ := NewBoardImp(b, config)
		if !imp.UpdateOnce() {
			t.Errorf("%v", b)
		}
		if imp.statics.combo != c.combo {
			t.Errorf("%v, expected: %v, actural: %v", b, c.combo, imp.statics.combo)
		}
		if String(imp.board.orbs) != c.updatedBoard {
			expected, _ := CopyBoard(b)
			expected.SetBoard(c.updatedBoard)
			t.Errorf("expected: %v, actural: %v", expected, imp.board.orbs)
		}
	}
}

func TestUpdate(t *testing.T) {
	b, _ := NewBoard(5, 6, "000000000000000000000000000000")
	config := &Config{3, 10}
	cases := []struct {
		board        string
		updatedBoard string
		combo        int
	}{
		{"566363644363422363233563111355", "000000000000000000000000500000", 9},
		{"245451243441233341243441222111", "000000000000000000000000005050", 6},
	}
	for _, c := range cases {
		b.SetBoard(c.board)
		imp, _ := NewBoardImp(b, config)
		imp.Update()
		if imp.statics.combo != c.combo {
			t.Errorf("%v, expected: %v, actural: %v", b, c.combo, imp.statics.combo)
		}
		if String(imp.board.orbs) != c.updatedBoard {
			expected, _ := CopyBoard(b)
			expected.SetBoard(c.updatedBoard)
			t.Errorf("expected: %v, actural: %v", expected, imp.board.orbs)
		}
	}
}

func TestFindPath(t *testing.T) {
	b, _ := NewBoard(5, 6, "000000000000000000000000000000")
	config := &Config{3, 10}
	// b.SetBoard("566363644363423631223563131355")
	b.SetBoard("263511332511644343336262246124")
	// b.SetBoard("111111111111111111111111111111")
	sol := FindPath(b, config)
	t.Logf("%v %v %v", sol.path, sol.statics, sol.board)
}
