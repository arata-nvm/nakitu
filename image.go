package main

import (
	"fmt"
	"io"
)

type Image struct {
	Width  int
	Height int
	Pixels []RGB
}

type RGB [3]int

func NewRGB(r, g, b int) RGB {
	return RGB{r, g, b}
}

func NewImage(width, height int) *Image {
	pixels := make([]RGB, width*height)
	return &Image{Width: width, Height: height, Pixels: pixels}
}

func (i *Image) SetPixel(x, y int, color RGB) {
	index := x + y*i.Width
	i.Pixels[index] = color
}

func (i *Image) Write(w io.Writer) {
	fmt.Fprintln(w, "P3")
	fmt.Fprintf(w, "%d %d\n", i.Width, i.Height)
	fmt.Fprintln(w, "255")

	for _, pixel := range i.Pixels {
		fmt.Fprintf(w, "%d %d %d\n", pixel[0], pixel[1], pixel[2])
	}
}
