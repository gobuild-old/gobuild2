package slave

import (
	"fmt"
	"runtime"
	"time"
)

type Mission struct {
	Repo   string
	Branch string
	Cgo    bool
}

var missionQueue = make(chan Mission)

func init() {
	n := runtime.NumCPU()
	n = 1
	for i := 0; i < n; i++ {
		go func() {
			for {
				var mission = <-missionQueue
				fmt.Println(mission)
				work(&mission)
				time.Sleep(time.Second)
			}
		}()
	}
}
