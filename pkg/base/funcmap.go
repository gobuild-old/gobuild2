package base

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"
	"time"

	"github.com/gobuild/gobuild2/pkg/units"
)

var TemplateFuncs = template.FuncMap{
	"title":     strings.Title,
	"ansi2html": ansi2html,
	"timesince": timeSince,
}

func ansi2html(s string) template.HTML {
	p := "\033[1;"
	h := `<span style="color:%s">`
	colorMap := map[string]string{
		"30": "#CCCCCC", //"gray",
		"31": "red",
		"32": "green",
		"33": "#FFFF00", //"yellow",
		"34": "#3366FF", //"blue",
		"35": "#9933CC", //"purple",
		"36": "#66CCFF",
	}
	for num, color := range colorMap {
		s = strings.Replace(s, p+num+"m", fmt.Sprintf(h, color), -1)
	}
	s = strings.Replace(s, "\033[0m", "</span>", -1)
	s = regexp.MustCompile("(\033"+`[;\[\d]+?m)+`).ReplaceAllString(s, `<span style="color:black">`)
	return template.HTML(s)
}

func timeSince(t time.Time) string {
	dur := time.Since(t) + time.Hour*8
	fmt.Println(t, dur+time.Hour*8)
	return units.HumanDuration(dur)
}
