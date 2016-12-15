package httpexporter

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Proxy struct {
	target      *url.URL
	proxy       *httputil.ReverseProxy
	metricsPort string
}

func New(target string, metricsPort string) *Proxy {
	url, _ := url.Parse(target)

	p := &Proxy{target: url, proxy: httputil.NewSingleHostReverseProxy(url)}
	p.metricsPort = metricsPort

	return p
}

func (p *Proxy) Handler() http.HandlerFunc {
	// Declare metrics
	var instLabels = []string{"method", "code"}
	reqCnt := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests",
			Help: "Number of http requests.",
		},
		instLabels,
	)
	prometheus.Register(reqCnt)

	reqDur := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "http_responsetime_microseconds",
			Help: "Summary about the response time of the http requests counted in microseconds",
		})
	prometheus.Register(reqDur)

	// Start dedicated prometheus webserver
	serverMuxProm := http.NewServeMux()
	serverMuxProm.Handle("/metrics", promhttp.Handler())
	serverMuxProm.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<html>
				<head><title>HTTP Exporter</title></head>
				<body>
					<h1>HTTP Exporter</h1>
					<p><a href="/metrics"">Metrics</a></p>
				</body>
			</html>`))
	})

	go func() {
		log.Fatal(http.ListenAndServe(p.metricsPort, serverMuxProm))
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		delegate := &responseWriterDelegator{ResponseWriter: w}

		var rw http.ResponseWriter
		rw = delegate

		rw.Header().Set("X-HTTPEXPORTER", "1")
		p.proxy.ServeHTTP(rw, r)

		elapsed := float64(time.Since(now)) / float64(time.Microsecond)

		// log.Printf("Status: %d Method: %s\n", delegate.status, strings.ToLower(r.Method))
		// log.Printf("Duration: %f\n", elapsed)

		reqDur.Observe(elapsed)
		reqCnt.WithLabelValues(strings.ToLower(r.Method), strconv.Itoa(delegate.status)).Inc()

	})

}
