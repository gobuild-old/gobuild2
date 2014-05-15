package slave

import (
	"fmt"
	"runtime"
	"time"

	"github.com/gobuild/gobuild2/pkg/xrpc"
)

var missionQueue = make(chan *xrpc.Mission)

func startWork() {
	n := runtime.NumCPU()
	n = 1
	for i := 0; i < n; i++ {
		go func() {
			for {
				var mission = <-missionQueue
				fmt.Println(mission)
				work(mission)
				time.Sleep(time.Second)
			}
		}()
	}
}
