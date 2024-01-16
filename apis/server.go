package apis

import (
	"go-boilerplate-api/apis/grpc"
	"go-boilerplate-api/apis/http"
	"go-boilerplate-api/shared"
	"sync"
)

// InitServers will pass the dependencies to the servers.
// The servers will start in an individual goroutine
// Wait group is used to wait for all the goroutines launched here to finish.
// In in ideal scenerio the routines would run indefinitely
// InitServers will pass the dependencies to the servers.
// The servers will start in an individual goroutine
// Wait group is used to wait for all the goroutines launched here to finish.
// In in ideal scenerio the routines would run indefinitely
func InitServers(deps *shared.Deps) error {
	var wg sync.WaitGroup
	var fatalErrChan = make(chan error)
	var wgDone = make(chan bool)

	wg.Add(2)
	go http.StartServer(deps, &wg, fatalErrChan)
	go grpc.StartServer(deps, &wg, fatalErrChan)

	// Final goroutine to wait until WaitGroup is done
	go func() {
		wg.Wait()
		close(wgDone)
	}()

	// Wait until either WaitGroup is done or an error is received through the channel
	select {
	case <-wgDone:
		// carry on
		break
		// Catch error from  channel and return
	case err := <-fatalErrChan:
		close(fatalErrChan)
		return err
	}

	return nil
}
