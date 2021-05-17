package monitoring

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type QueryQueued struct {
	PoolName                string `db:"pool_name"`
	QueuedQueriesCount      int    `db:"queued_queries_count"`
	QueuedQueriesDurationS  int    `db:"queued_queries_duration_second"`
}

func NewQueryQueued(db *sqlx.DB) []QueryQueued {
	sql := `
	SELECT
    pool_name,
    COUNT(transaction_id) as queued_queries_count,
    MAX(TIMESTAMPDIFF(S, queue_entry_timestamp, CLOCK_TIMESTAMP()))::INT as queued_queries_duration_second
  FROM v_monitor.RESOURCE_QUEUES
  GROUP BY pool_name;
	`

	queryQueued := []QueryQueued{}
	err := db.Select(&queryQueued, sql)
	if err != nil {
		log.Fatal(err)
	}

	return queryQueued
}

func (qr QueryQueued) ToMetric() map[string]int {
	metrics := map[string]int{}

	poolname := fmt.Sprintf("pool_name=%q", qr.PoolName)
	metrics[fmt.Sprintf("vertica_queued_request_duration_s{%s}", poolname)] = qr.QueuedQueriesDurationS
	metrics[fmt.Sprintf("vertica_queued_request_count{%s}", poolname)] = qr.QueuedQueriesCount

	return metrics
}