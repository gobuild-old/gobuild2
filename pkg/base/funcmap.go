package base

import (
	"fmt"
	"html/template"
	"strings"
)

var TemplateFuncs = template.FuncMap{
	"title":     strings.Title,
	"ansi2html": ansi2html,
}

func ansi2html(s string) string {
	p := "\033[1;"
	h := "<span color='%s'>"
	colorMap := map[string]string{
		"30": "#CCCCCC", //"gray",
		"31": "red",
		"32": "green",
		"33": "yellow",
		"34": "#3366FF", //"blue",
		"35": "#9933CC", //"purple",
		"36": "#66CCFF",
	}
	for num, color := range colorMap {
		s = strings.Replace(s, p+num+"m", fmt.Sprintf(h, color), -1)
	}
	return s
}
