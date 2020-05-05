package utcpq

import (
	"database/sql"

	"github.com/hirofumi/utcer"
	"github.com/lib/pq"
)

// nolint: gochecknoinits
func init() {
	sql.Register("utcpq", utcer.Wrap(&pq.Driver{}))
}
