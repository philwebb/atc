package api

import (
	"net/http"

	"github.com/concourse/atc/metric"
)

type CTServer struct{}

func (s *CTServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ct := make(map[*int][]byte)
	metric.CTMutex.Lock()
	for k, v := range metric.ConnectionTraces {
		ct[k] = v
	}
	metric.CTMutex.Unlock()

	for _, v := range ct {
		w.Write(v)
		w.Write([]byte("\n\n\n"))
	}
}
