package main

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"

	_ "embed"
)

//go:embed "template.png"
var rawImage []byte

//go:embed "index.html"
var rawIndex []byte

type App struct {
	ff font.Face
	bg image.Image
}

func NewApp() (ret App) {
	// Font
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}
	ret.ff = truetype.NewFace(font, &truetype.Options{Size: 48})

	// Background
	ret.bg, err = png.Decode(bytes.NewReader(rawImage))
	if err != nil {
		log.Fatal(err)
	}

	return
}

func (app *App) render(w io.Writer, lineOne string, lineTwo string) {
	const (
		W = 733
		H = 906
	)

	dc := gg.NewContext(W, H)
	dc.SetFontFace(app.ff)
	dc.SetRGB(0, 0, 0)

	dc.DrawImage(app.bg, 0, 0)
	dc.DrawStringWrapped(lineOne, 558, 114, 0.5, 0.5, 300, 1.1, gg.AlignCenter)
	dc.DrawStringWrapped(lineTwo, 571, 403, 0.5, 0.5, 300, 1.1, gg.AlignCenter)
	if err := dc.EncodePNG(w); err != nil {
		log.Fatal(err)
	}
}

func (app *App) index(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s from %s\n", req.Method, req.RequestURI, req.RemoteAddr)


	lone := req.FormValue("l1")
	ltwo := req.FormValue("l2")

	if lone == "" && ltwo == "" {
		w.Header().Set("Content-Type", "text/html")
		w.Write(rawIndex)
	} else {
		w.Header().Set("Content-Type", "image/png")
		app.render(w, lone, ltwo)
	}
}

func main() {
	app := NewApp()
	http.HandleFunc("/", app.index)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
