package main

import (
	"bytes"
	"flag"
	"image"
	"image/png"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/marpaia/graphite-golang"
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
	g  *graphite.Graphite
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
	start := time.Now()
	lone := req.FormValue("l1")
	ltwo := req.FormValue("l2")

	if lone == "" && ltwo == "" {
		w.Header().Set("Content-Type", "text/html")
		w.Write(rawIndex)
		app.g.SimpleSend("meme.index", strconv.FormatInt(time.Now().Sub(start).Milliseconds(), 10))
	} else {
		w.Header().Set("Content-Type", "image/png")
		app.render(w, lone, ltwo)
		app.g.SimpleSend("meme.render", strconv.FormatInt(time.Now().Sub(start).Milliseconds(), 10))
	}
}

func main() {
	cliGraphite := flag.String("graphite", "", "Remote Graphite host:port")
	flag.Parse()

	app := NewApp()

	if *cliGraphite == "" {
		app.g = graphite.NewGraphiteNop("", 0)
	} else {
		// host:prot splitter
		host, sport, err := net.SplitHostPort(*cliGraphite)
		if err != nil {
			log.Fatal(err)
		}
		// number conv
		port, err := strconv.Atoi(sport)
		if err != nil {
			log.Fatal(err)
		}
		// Set prefix
		app.g, err = graphite.NewGraphiteWithMetricPrefix(host, port, "maas")
		if err != nil {
			log.Fatal(err)
		}
	}

	http.HandleFunc("/", app.index)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
