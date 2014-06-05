package slave

import (
	"time"

	"github.com/gobuild/gobuild2/pkg/xrpc"
)

var missionQueue = make(chan *xrpc.Mission)

func startWork() {
	// n := runtime.NumCPU()
	n := 10
	for i := 0; i < n; i++ {
		go func() {
			for {
				var mission = <-missionQueue
				work(mission)
				time.Sleep(time.Second)
			}
		}()
	}
}
