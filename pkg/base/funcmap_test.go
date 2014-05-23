package base

import (
	"fmt"
	"strconv"
	"testing"
)

func TestAnsi2html(t *testing.T) {
	start := "\033[1;"
	clear := "\033[0m"
	for i := 30; i < 40; i++ {
		fmt.Println(start + strconv.Itoa(i) + "m" + "i see you:" + strconv.Itoa(i) + clear)
	}
}
