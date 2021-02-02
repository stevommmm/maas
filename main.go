package main

import (
	"log"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

func main() {
	const (
		W = 733
		H = 906
	)

	// Font
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}
	face := truetype.NewFace(font, &truetype.Options{Size: 48})

	// Background
	im, err := gg.LoadPNG("template.png")
	if err != nil {
		log.Fatal(err)
	}

	dc := gg.NewContext(W, H)
	dc.SetFontFace(face)
	dc.SetRGB(0, 0, 0)

	dc.DrawImage(im, 0, 0)
	dc.DrawStringWrapped("Hello, world! and other stuff", 558, 114, 0.5, 0.5, 300, 1.1, gg.AlignCenter)
	dc.DrawStringWrapped("Meow?", 571, 403, 0.5, 0.5, 300, 1.1, gg.AlignCenter)
	dc.SavePNG("out.png")
}
