package monitoring

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type QueryFailed struct {
	UserName                string `db:"user_name"`
	FailedQueriesCount      int    `db:"failed_queries_count"`
}

// Failed queries for last 10 minutes
func NewQueryFailed(db *sqlx.DB) []QueryFailed {
	sql := `
	SELECT
    user_name,
    COUNT(REQUEST_ID) as failed_queries_count
  FROM v_monitor.query_requests
  WHERE success = 'false' AND is_executing = 'false' AND TIMESTAMPDIFF(mi, start_timestamp, CLOCK_TIMESTAMP()) < 10
  GROUP BY user_name;
	`

	queryFailed := []QueryFailed{}
	err := db.Select(&queryFailed, sql)
	if err != nil {
		log.Fatal(err)
	}

	return queryFailed
}

func (qr QueryFailed) ToMetric() map[string]int {
	metrics := map[string]int{}

	username := fmt.Sprintf("user_name=%q", qr.UserName)
	metrics[fmt.Sprintf("vertica_failed_request_count{%s}", username)] = qr.FailedQueriesCount

	return metrics
}