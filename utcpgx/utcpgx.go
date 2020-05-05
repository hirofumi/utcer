package utcpgx

import (
	"database/sql"

	"github.com/hirofumi/utcer"
	"github.com/jackc/pgx/v4/stdlib"
)

// nolint: gochecknoinits
func init() {
	sql.Register("utcpgx", utcer.Wrap(stdlib.GetDefaultDriver()))
}
