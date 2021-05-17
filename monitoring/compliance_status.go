package monitoring

import (
	"log"
	"strings"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type ComplianceStatusQuery struct {
	ComplianceStatus          string `db:"compliance_status"`
}

func NewComplianceStatus(db *sqlx.DB) []ComplianceStatusQuery {
	sql := `
	SELECT GET_COMPLIANCE_STATUS() as compliance_status;
	`
	complianceStatus := []ComplianceStatusQuery{}
	err := db.Select(&complianceStatus, sql)

	if err != nil {
		log.Fatal(err)
	}

	return complianceStatus
}

// ToMetric converts ComplianceStatus to a Map.
func (qr ComplianceStatusQuery) ToMetric() map[string]int {
// Compliance utilization response looks like:
// --------------
// Raw Data Size: 0.69TB +/- 0.07TB
//  License Size : 1.00TB
//  Utilization  : 69%
//  Audit Time   : 2021-05-10 23:59:11.051827+00
//  Node count : 3
//  License Node limit : 3
//  Compliance Status : The database is in compliance with respect to raw data size.
//
//  No expiration date for a Perpetual license
// --------------

	metrics := map[string]int{}
  listCompliance := strings.Split(qr.ComplianceStatus, "\n")
  percent, err := strconv.Atoi(strings.TrimSpace(strings.TrimSuffix(strings.Split(listCompliance[2], ":")[1], "%")))
  if err != nil {
  		log.Fatal(err)
  }
	metrics["vertica_compliance_utilization_percentage"] = percent

	return metrics
}
