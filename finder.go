package pazudoraer

import "sort"

// Up, Down, Left, Right
var MOVES = [4]Position{Position{-1, 0}, Position{1, 0}, Position{0, -1}, Position{0, 1}}

type Solution struct {
	board   *Board
	path    []Position
	statics *Statics
	point   int
	isMove  bool
}

func NewSolution(b *Board) (*Solution, error) {
	s := new(Solution)
	s.board, _ = CopyBoard(b)
	s.path = make([]Position, 0)
	s.statics = &Statics{}
	return s, nil
}

func CopySolution(sol *Solution) (*Solution, error) {
	s := new(Solution)
	s.board, _ = CopyBoard(sol.board)
	s.path = make([]Position, 0, len(sol.path)+1)
	for _, val := range sol.path {
		s.path = append(s.path, val)
	}
	return s, nil
}

type byPoint []*Solution

func (s byPoint) Len() int {
	return len(s)
}
func (s byPoint) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byPoint) Less(i, j int) bool {
	return s[i].point < s[j].point
}

type Ranker interface {
	CalcPoint(s *Solution) int
}

type ComboPoint struct{}

func (c *ComboPoint) CalcPoint(s *Solution) int {
	return s.statics.combo
}

type Finder interface {
	FindPath(b *Board, c *Config) *Solution
}

type PrunedFinder struct {
	ranker  Ranker
	perStep int
	keepNum int
}

type FullFinder struct{}

func (p *PrunedFinder) FindPath(b *Board, c *Config) *Solution {
	height, width := b.shape.height, b.shape.width
	sols := make([]*Solution, 0, 4*height*width)

	// set all points and first move as starting point
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			currPosition := Position{i, j}
			for _, move := range MOVES {
				nextPosition := Add(currPosition, move)

				// swap two orbs and add to solutions
				if b.InBoard(nextPosition) {
					newSol, _ := NewSolution(b)
					newSol.board.Swap(currPosition, nextPosition)
					// skip the board same as initial board
					if String(newSol.board.orbs) == String(b.orbs) {
						continue
					}
					newSol.path = append(newSol.path, currPosition, nextPosition)
					sols = append(sols, newSol)
				}
			}
		}
	}

	for t := 0; t < c.maxPathLength/p.perStep; t++ {
		// move and prune
		start, end := 0, len(sols)
		for i := 0; i < p.perStep; i++ {
			for j := start; j < end; j++ {
				sol := sols[j]
				if sol.isMove {
					continue
				}
				currPosition := sol.path[len(sol.path)-1]
				for _, move := range MOVES {
					nextPosition := Add(currPosition, move)

					// skip go back move
					prevPosition := sol.path[len(sol.path)-2]
					if prevPosition == nextPosition {
						continue
					}

					// swap two orbs and add to solutions
					if sol.board.InBoard(nextPosition) {
						newSol, _ := CopySolution(sol)
						newSol.board.Swap(currPosition, nextPosition)
						newSol.path = append(newSol.path, nextPosition)
						sols = append(sols, newSol)
					}
				}
				sol.isMove = true
			}
			start, end = end, len(sols)
		}

		// evaluate and prune
		for _, s := range sols {
			if s.statics == nil {
				b, _ := CopyBoard(s.board)
				s.statics, _ = Update(b, c)
				s.point = p.ranker.CalcPoint(s)
			}
		}
		sort.Sort(sort.Reverse(byPoint(sols)))
		sols = sols[:p.keepNum]
	}

	return sols[0]
}

func (p *FullFinder) FindPath(b *Board, c *Config) *Solution {
	height, width := b.shape.height, b.shape.width
	sols := make([]*Solution, 0, 4*height*width)
	bestSol, _ := NewSolution(b)
	bestSol.board.Swap(Position{0, 0}, Position{0, 1})
	bestSol.path = append(bestSol.path, Position{0, 0}, Position{0, 1})
	bestSol.statics, _ = Update(bestSol.board, c)

	// set all points and first move as starting point
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			currPosition := Position{i, j}
			for _, move := range MOVES {
				nextPosition := Add(currPosition, move)

				// swap two orbs and add to solutions
				if b.InBoard(nextPosition) {
					newSol, _ := NewSolution(b)
					newSol.board.Swap(currPosition, nextPosition)
					// skip the board same as initial board
					if String(newSol.board.orbs) == String(b.orbs) {
						continue
					}
					newSol.path = append(newSol.path, currPosition, nextPosition)
					sols = append(sols, newSol)
				}
			}
		}
	}

	for len(sols) > 0 {
		sol := sols[0]
		sols = sols[1:]

		if len(sol.path) < c.maxPathLength {
			currPosition := sol.path[len(sol.path)-1]
			for _, move := range MOVES {
				nextPosition := Add(currPosition, move)

				// skip go back move
				prevPosition := sol.path[len(sol.path)-2]
				if prevPosition == nextPosition {
					continue
				}

				// swap two orbs and add to solutions
				if sol.board.InBoard(nextPosition) {
					newSol, _ := CopySolution(sol)
					newSol.board.Swap(currPosition, nextPosition)
					newSol.path = append(newSol.path, nextPosition)
					sols = append(sols, newSol)
				}
			}
		}

		// evaluate solution
		sol.statics, _ = Update(sol.board, c)
		if sol.statics.combo > bestSol.statics.combo {
			bestSol = sol
		} else if sol.statics.combo == bestSol.statics.combo && len(sol.path) < len(bestSol.path) {
			bestSol = sol
		}
	}

	return bestSol
}
