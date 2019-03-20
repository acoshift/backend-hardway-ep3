package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var reg = prometheus.NewRegistry()

// docker run --name prom -d -p 9090:9090 -v $(pwd)/prom/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
// getent hosts host.docker.internal | cut -d' ' -f1
// docker run --name grafana -d -p 3000:3000 grafana/grafana
// grafana default user/password = admin/admin
func main() {
	reg.MustRegister(prometheus.NewGoCollector())
	reg.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{
		PidFn: func() (int, error) { return os.Getpid(), nil },
	}))

	go http.ListenAndServe(":9000", promhttp.InstrumentMetricHandler(
		reg,
		promhttp.HandlerFor(reg, promhttp.HandlerOpts{
			DisableCompression: true,
		}),
	))

	h := http.HandlerFunc(handler)
	h = requestTracker(h)
	http.ListenAndServe(":8080", h)
}

func requestTracker(h http.Handler) http.HandlerFunc {
	requests := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "server",
		Name:      "requests",
	}, []string{"host", "status", "method"})
	reg.MustRegister(requests)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := prometheus.Labels{
			"method": r.Method,
			"host":   r.Host,
		}
		nw := requestTrackRW{
			ResponseWriter: w,
		}
		defer func() {
			l["status"] = strconv.Itoa(nw.status)
			counter, err := requests.GetMetricWith(l)
			if err != nil {
				return
			}
			counter.Inc()
		}()

		h.ServeHTTP(&nw, r)
	})
}

type requestTrackRW struct {
	http.ResponseWriter

	wroteHeader bool
	status      int
}

func (w *requestTrackRW) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *requestTrackRW) Write(p []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(p)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
