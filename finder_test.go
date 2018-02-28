package pazudoraer

import "testing"

func TestFindPath(t *testing.T) {
	b, _ := NewBoard(5, 6, "000000000000000000000000000000")
	config := &Config{3, 40}
	// b.SetBoard("566363644363423631223563131355")
	// b.SetBoard("263511332511644343336262246124")
	// b.SetBoard("111111111111111111111111111111")
	// b.SetBoard("121212212121121212212121121212")
	b.SetBoard("632131115436456252643451452566")
	f := PrunedFinder{&ComboPoint{}, 4, 256}
	sol := f.FindPath(b, config)
	t.Errorf("%v %#v %v", sol.path, sol.statics, sol.board)
}
