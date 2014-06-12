package routers

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"

	"code.google.com/p/freetype-go/freetype"
)

const (
	fontFile = "public/fonts/font.ttf"
	fontSize = 12
	fontDPI  = 72
)

var (
	black color.Color = color.RGBA{50, 50, 50, 255}
	green color.Color = color.RGBA{85, 154, 17, 200}
)

func Badge(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "image/png")
	BadgeEncode(w, "download")
}

func BadgeEncode(w io.Writer, message string) (err error) {
	// draw "gobuild.io | download"
	webname := "GOBUILD"
	const middle = 65
	const gap = 8
	img := image.NewNRGBA(image.Rect(0, 0, middle+(len(message)+1)*7+gap*2, 18))
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		return
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return
	}
	left, right := img.Bounds(), img.Bounds()
	left.Max = image.Pt(middle, 18)
	right.Min = image.Pt(middle, 0)
	// fill left(black) right(green)
	draw.Draw(img, left, &image.Uniform{black}, image.ZP, draw.Src)
	draw.Draw(img, right, &image.Uniform{green}, image.ZP, draw.Src)

	c := freetype.NewContext()
	c.SetDPI(fontDPI)
	c.SetFont(font)
	c.SetFontSize(fontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.White)
	pt := freetype.Pt(gap, 12)
	_, err = c.DrawString(webname, pt) // 10 chars width = 60px
	if err != nil {
		return
	}
	c.SetSrc(image.Black)
	pt = freetype.Pt(middle+gap, 12)
	_, err = c.DrawString(message, pt)

	// w.Header().Set("Content-Type", "image/png")
	return png.Encode(w, img)
}
