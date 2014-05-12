package slave

import (
	"fmt"
	"runtime"
)

type Mission struct {
	Repo   string
	Branch string
	Cgo    bool
}

var missionQueue = make(chan Mission)

func init() {
	n := runtime.NumCPU()
	for i := 0; i < n; i++ {
		go func() {
			var mission = <-missionQueue
			fmt.Println(mission)
			work(&mission)
		}()
	}
}
