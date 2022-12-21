package internal

import (
	"encoding/json"
	"fmt"
)

type Canvas []int

func NewCanvas() *Canvas {
	canvas := make([]int, 10000)
	for i := 0; i < 10000; i++ {
		canvas[i] = 0xffffff
	}
	var output Canvas = canvas
	return &output
}

func (c Canvas) MarshalJSON() ([]byte, error) {
	var canvas []int = c

	output := struct {
		Type   string
		Canvas []int
	}{
		Type:   "draw",
		Canvas: canvas,
	}

	return json.Marshal(output)
}

func (c Canvas) SetCoordinate(x, y int, color int) {
	c[100*y+x] = color
}

func (c Canvas) String() string {
	s := ""
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			if c[100*i+j] != 0xffffff {
				s += fmt.Sprint(i) + " " + fmt.Sprint(j) + " " + fmt.Sprint(c[100*i+j]) + "\n"
			}
		}
	}
	return s
}
