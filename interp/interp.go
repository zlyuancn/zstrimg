package interp

import (
	"image"
	"image/color"
)

type Interp interface {
	Interp(src image.Image, x, y float64) color.Color
}

type RGBA interface {
	RGBA(src *image.RGBA, x, y float64) color.RGBA
}

type Gray interface {
	Gray(src *image.Gray, x, y float64) color.Gray
}
