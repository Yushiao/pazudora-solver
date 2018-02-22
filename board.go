package pazudoraer

import (
	"fmt"
	"strconv"
	"strings"
)

type Orb int

const (
	NULL Orb = iota
	FIRE
	WATER
	WOOD
	LIGHT
	DARK
	HEART
	JAMMER
	POISON
	MORTAL_POISON
)

type Position struct {
	height, width int
}

type Board struct {
	orbs  []Orb
	shape Position
}

type Statics struct {
	maxCombo  int
	combo     int
	fallCount int
}

type BoardImp struct {
	board   *Board
	config  *Config
	statics *Statics

	matched  []bool
	checked  []bool
	location []Position
	orbs     []Orb
}

type Config struct {
	minMatchNum   int
	maxPathLength int
}

// Up, Down, Left, Right
var MOVES = [4]Position{Position{-1, 0}, Position{1, 0}, Position{0, -1}, Position{0, 1}}

type Solution struct {
	board   *Board
	path    []Position
	statics *Statics
}

// Set board by string of numbers.
// like "111222333111111222333111111222"
func NewBoard(height, width int, board string) (*Board, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("width or height smaller than zero")
	}
	b := &Board{make([]Orb, 0, height*width), Position{height, width}}
	err := b.SetBoard(board)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (b *Board) SetBoard(board string) error {
	height, width := b.shape.height, b.shape.width
	if len(board) != width*height {
		return fmt.Errorf("size not match")
	}
	b.orbs = b.orbs[:0]
	for _, orb := range strings.Split(board, "") {
		i, _ := strconv.Atoi(orb)
		b.orbs = append(b.orbs, Orb(i))
	}
	return nil
}

func CopyBoard(b *Board) (*Board, error) {
	return NewBoard(b.shape.height, b.shape.width, String(b.orbs))
}

func Update(b *Board, c *Config) (*Statics, *Board) {
	imp, _ := NewBoardImp(b, c)
	imp.Update()
	return imp.statics, imp.board
}

func (b *Board) Swap(src, dst Position) {
	srcIndex := src.height*b.shape.width + src.width
	dstIndex := dst.height*b.shape.width + dst.width
	b.orbs[srcIndex], b.orbs[dstIndex] = b.orbs[dstIndex], b.orbs[srcIndex]
}

func (b *Board) InBoard(p Position) bool {
	height, width := b.shape.height, b.shape.width
	return p.height >= 0 && p.height < height && p.width >= 0 && p.width < width
}

func NewBoardImp(b *Board, c *Config) (*BoardImp, error) {
	imp := new(BoardImp)
	imp.board = b
	imp.config = c
	imp.statics = &Statics{}

	height, width := b.shape.height, b.shape.width
	imp.matched = make([]bool, height*width)
	imp.checked = make([]bool, height*width)
	imp.location = make([]Position, 0, height*width)
	imp.orbs = make([]Orb, 0, height)
	return imp, nil
}

func (b *BoardImp) Update() {
	for b.UpdateOnce() {
		b.statics.fallCount += 1
	}
}

func (b *BoardImp) UpdateOnce() bool {
	// clear values
	for index := range b.matched {
		b.matched[index] = false
		b.checked[index] = false
	}
	b.location = b.location[:0]
	b.orbs = b.orbs[:0]

	isMatch := b.checkMatch()
	if !isMatch {
		return false
	}

	b.updateStatics()
	b.floodFill()
	return true
}

// Check orbs are match
func (b *BoardImp) checkMatch() bool {
	isMatch := false
	height, width := b.board.shape.height, b.board.shape.width
	minMatchNum := b.config.minMatchNum

	// match horizontal
	for i := 0; i < height; i++ {
		for j := 0; j < width-minMatchNum+1; j++ {
			if b.checkMatchHor(i, j) {
				isMatch = true
				for k := 0; k < minMatchNum; k++ {
					now := i*width + j + k
					b.matched[now] = true
				}
			}
		}
	}

	// match vertical
	for i := 0; i < height-minMatchNum+1; i++ {
		for j := 0; j < width; j++ {
			if b.checkMatchVer(i, j) {
				isMatch = true
				for k := 0; k < minMatchNum; k++ {
					now := (i+k)*width + j
					b.matched[now] = true
				}
			}
		}
	}

	return isMatch
}

// Check orbs are match horizontally from a certain starting point
func (b *BoardImp) checkMatchHor(h, w int) bool {
	width := b.board.shape.width
	minMatchNum := b.config.minMatchNum
	start := h*width + w
	if b.board.orbs[start] == NULL {
		return false
	}
	for i := 1; i < minMatchNum; i++ {
		now := h*width + w + i
		if b.board.orbs[now] != b.board.orbs[start] {
			return false
		}
	}
	return true
}

// Check orbs are match vertically from a certain starting point
func (b *BoardImp) checkMatchVer(h, w int) bool {
	width := b.board.shape.width
	minMatchNum := b.config.minMatchNum
	start := h*width + w
	if b.board.orbs[start] == NULL {
		return false
	}
	for i := 1; i < minMatchNum; i++ {
		now := (h+i)*width + w
		if b.board.orbs[now] != b.board.orbs[start] {
			return false
		}
	}
	return true
}

// Update board statics
func (b *BoardImp) updateStatics() {
	height, width := b.board.shape.height, b.board.shape.width
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			start := i*width + j
			if b.checked[start] {
				continue
			}
			if !b.matched[start] {
				b.checked[start] = true
				continue
			}
			b.statics.combo += 1
			b.location = append(b.location, Position{i, j})
			for len(b.location) > 0 {
				pos := b.location[0]
				b.location = b.location[1:]
				index := pos.height*width + pos.width
				if b.matched[index] && !b.checked[index] && b.board.orbs[start] == b.board.orbs[index] {
					b.checked[index] = true
					// check previous row
					if pos.height-1 >= 0 {
						b.location = append(b.location, Position{pos.height - 1, pos.width})
					}
					// check next row
					if pos.height+1 < height {
						b.location = append(b.location, Position{pos.height + 1, pos.width})
					}
					// check previous column
					if pos.width-1 >= 0 {
						b.location = append(b.location, Position{pos.height, pos.width - 1})
					}
					// check next column
					if pos.width+1 < width {
						b.location = append(b.location, Position{pos.height, pos.width + 1})
					}
				}
			}
		}
	}
}

// Flood fill orbs
func (b *BoardImp) floodFill() {
	height, width := b.board.shape.height, b.board.shape.width
	for j := 0; j < width; j++ {
		b.orbs = b.orbs[:0]
		for i := 0; i < height; i++ {
			index := (height-i-1)*width + j
			if !b.matched[index] && b.board.orbs[index] != NULL {
				b.orbs = append(b.orbs, b.board.orbs[index])
			}
		}
		needed := height - len(b.orbs)
		for i := 0; i < needed; i++ {
			b.orbs = append(b.orbs, NULL)
		}
		for i := 0; i < height; i++ {
			index := (height-i-1)*width + j
			b.board.orbs[index] = b.orbs[i]
		}
	}
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

func FindPath(b *Board, c *Config) *Solution {
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

// Add two position
func Add(p, q Position) Position {
	return Position{p.height + q.height, p.width + q.width}
}

func NewConfig() *Config {
	return &Config{3, 30}
}

func (s *Solution) Get() (*Board, []Position, *Statics) {
	return s.board, s.path, s.statics
}
