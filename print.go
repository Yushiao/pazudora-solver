package pazudoraer

import (
	"bytes"
	"fmt"
)

func String(a []Orb) string {
	var buffer bytes.Buffer
	for i := range a {
		buffer.WriteString(fmt.Sprintf("%d", a[i]))
	}
	return buffer.String()
}

func (b *Board) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("(%d, %d)\n", b.shape.height, b.shape.width))
	for i := 0; i < b.shape.height; i++ {
		start, end := i*b.shape.width, (i+1)*b.shape.width
		buffer.WriteString(fmt.Sprintf("%v\n", b.orbs[start:end]))
	}
	return buffer.String()
}

func (s *Solution) String() string {
	var buffer bytes.Buffer
	for _, val := range s.board.orbs {
		buffer.WriteString(fmt.Sprintf("%d", val))
	}
	pos := s.path[len(s.path)-1]
	buffer.WriteString(fmt.Sprintf("%d%d", pos.height, pos.width))
	return buffer.String()
}
