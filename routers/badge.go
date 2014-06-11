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
	BadgeEncode(w)
}

func BadgeEncode(w io.Writer) (err error) {
	img := image.NewNRGBA(image.Rect(0, 0, 140, 18))
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		return
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return
	}
	left, right := img.Bounds(), img.Bounds()
	const middle = 65
	left.Max = image.Pt(middle, 18)
	right.Min = image.Pt(middle, 0)
	// fill left(black) right(green)
	draw.Draw(img, left, &image.Uniform{black}, image.ZP, draw.Src)
	draw.Draw(img, right, &image.Uniform{green}, image.ZP, draw.Src)

	// draw "gobuild.io | download"
	c := freetype.NewContext()
	c.SetDPI(fontDPI)
	c.SetFont(font)
	c.SetFontSize(fontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.White)
	pt := freetype.Pt(7, 12)
	_, err = c.DrawString("gobuild.io", pt) // 10 chars width = 60px
	if err != nil {
		return
	}
	c.SetSrc(image.Black)
	pt = freetype.Pt(middle+13, 12)
	_, err = c.DrawString("download", pt)

	// w.Header().Set("Content-Type", "image/png")
	return png.Encode(w, img)
}
