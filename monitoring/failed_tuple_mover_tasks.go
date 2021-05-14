package monitoring

import (
	"log"
	"github.com/jmoiron/sqlx"
)

type FailedTupleMoverTasks struct {
	Count          int `db:"failed_tasks"`
}

func NewFailedTupleMoverTasks(db *sqlx.DB) []FailedTupleMoverTasks {
	sql := `
	SELECT COALESCE(count(*),0) as failed_tasks
  from tuple_mover_operations
  where operation_status = 'Abort';
	`
	failedTupleMoverTasks := []FailedTupleMoverTasks{}
	err := db.Select(&failedTupleMoverTasks, sql)

	if err != nil {
		log.Fatal(err)
	}

	return failedTupleMoverTasks
}

// ToMetric converts FailedTupleMoverTasks to a Map.
func (qr FailedTupleMoverTasks) ToMetric() map[string]int {

	metrics := map[string]int{}
	metrics["vertica_failed_tuple_mover_operations"] = qr.Count

	return metrics
}
