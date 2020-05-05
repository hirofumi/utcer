package pqutc

import (
	"database/sql"

	"github.com/hirofumi/utcer"
	"github.com/lib/pq"
)

// nolint: gochecknoinits
func init() {
	sql.Register("pqutc", utcer.Wrap(&pq.Driver{}))
}
