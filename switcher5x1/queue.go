package switcher5x1

import (
	"net/http"
	"sync"
	"time"

	"github.com/byuoitav/common/log"
)

type switcher5x1Response struct {
	response *http.Response
	err      error
}

type switcher5x1Request struct {
	request  *http.Request
	complete chan switcher5x1Response
}

var (
	queue                   = make(chan switcher5x1Request)
	once                    = sync.Once{}
	listeningFor5x1Requests = false
)

func processRequests() {
	for req := range queue {
		resp, error := http.DefaultClient.Do(req.request)
		req.complete <- switcher5x1Response{response: resp, err: error}
		//pause a little
		time.Sleep(250)
	}
}

func make5x1Request(req *http.Request) (*http.Response, error) {
	once.Do(func() {
		log.L.Infof("Launching request listener")
		go processRequests()
	})

	responseChannel := make(chan switcher5x1Response)
	queue <- switcher5x1Request{request: req, complete: responseChannel}
	response := <-responseChannel
	return response.response, response.err
}
