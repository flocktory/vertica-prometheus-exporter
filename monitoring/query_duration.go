package monitoring

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type QueryDuration struct {
	UserName                string `db:"user_name"`
	RunningRequestDurationS int    `db:"running_request_duration_second"`
}

// NewQueryRequests returns query performance for all users.
func NewQueryDuration(db *sqlx.DB) []QueryDuration {
	sql := `
	SELECT
    user_name,
    MAX(TIMESTAMPDIFF(S, start_timestamp, CLOCK_TIMESTAMP())) as running_request_duration_second
    FROM v_monitor.query_requests
    WHERE is_executing = 'true'
    GROUP BY user_name;
	`

	queryDuration := []QueryDuration{}
	err := db.Select(&queryDuration, sql)
	if err != nil {
		log.Fatal(err)
	}

	return queryDuration
}

// ToMetric converts QueryRequest to a Map.
func (qr QueryDuration) ToMetric() map[string]int {
	metrics := map[string]int{}

	username := fmt.Sprintf("user_name=%q", qr.UserName)
	metrics[fmt.Sprintf("vertica_running_request_duration_s{%s}", username)] = qr.RunningRequestDurationS

	return metrics
}