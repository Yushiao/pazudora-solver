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

// Add two position
func Add(p, q Position) Position {
	return Position{p.height + q.height, p.width + q.width}
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

// NewBoard set board by string of numbers.
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
		b.statics.fallCount++
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
			b.statics.combo++
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
