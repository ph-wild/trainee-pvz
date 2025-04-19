package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const AppName = "pvz_service"

const (
	labelApp    = "app"
	labelPath   = "path"
	labelCode   = "code"
	labelMethod = "method"
	labelEntity = "entity"
)

type Metrics struct {
	httpDurationSummary *prometheus.SummaryVec
	entityCount         *prometheus.CounterVec
}

func InitMetrics() *Metrics {
	m := &Metrics{}
	m.httpDurationSummary = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "http_request_duration_summary",
		Help:       "Duration of HTTP requests Summary.",
		Objectives: map[float64]float64{0.5: 0.5, 0.9: 0.9, 1: 1},
		AgeBuckets: 3,
		MaxAge:     120 * time.Second,
	}, []string{labelApp, labelPath, labelCode, labelMethod})
	prometheus.Register(m.httpDurationSummary)

	m.entityCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "created_entity_count",
		Help: "Count of created business entities.",
	}, []string{labelApp, labelEntity})
	prometheus.Register(m.entityCount)

	return m
}

func (m *Metrics) SaveHTTPDuration(timeSince time.Time, path string, code int, method string) {
	m.httpDurationSummary.With(map[string]string{
		labelApp:    AppName,
		labelPath:   path,
		labelCode:   strconv.Itoa(code),
		labelMethod: method,
	}).Observe(float64(time.Since(timeSince).Seconds()))
}

func (m *Metrics) SaveEntityCount(value float64, entity string) {
	m.entityCount.With(map[string]string{
		labelApp:    AppName,
		labelEntity: entity,
	}).Add(value)
}
