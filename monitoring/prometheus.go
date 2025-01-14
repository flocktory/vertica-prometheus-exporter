package monitoring

import (
	"github.com/jmoiron/sqlx"
	"regexp"
	"strings"
)

// PrometheusMetric maps a struct to a Prometheus valid map.
type PrometheusMetric interface {
	ToMetric() map[string]int
}

func NewPrometheusMetrics(db sqlx.DB) []PrometheusMetric {
	var metrics []PrometheusMetric

	for _, state := range NewNodeState(&db) {
		metrics = append(metrics, state)
	}
	for _, rejection := range NewPoolRejections(&db) {
		metrics = append(metrics, rejection)
	}
	for _, queryRequest := range NewQueryRequests(&db) {
		metrics = append(metrics, queryRequest)
	}
	for _, usage := range NewPoolUsage(&db) {
		metrics = append(metrics, usage)
	}
	for _, compliance := range NewComplianceStatus(&db) {
  		metrics = append(metrics, compliance)
  }
  for _, failedTasks := range NewFailedTupleMoverTasks(&db) {
    		metrics = append(metrics, failedTasks)
  }
  for _, queryDuration := range NewQueryDuration(&db) {
      		metrics = append(metrics, queryDuration)
  }
  for _, queryQueued := range NewQueryQueued(&db) {
        		metrics = append(metrics, queryQueued)
  }
  for _, queryFailed := range NewQueryFailed(&db) {
          		metrics = append(metrics, queryFailed)
  }
	metrics = append(metrics, NewVerticaSystem(&db))

	return metrics
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// ToSnakeCase converts all string values to snake case.
func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
