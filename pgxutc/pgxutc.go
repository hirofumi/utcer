package pgxutc

import (
	"database/sql"

	"github.com/hirofumi/utcer"
	"github.com/jackc/pgx/v4/stdlib"
)

func init() {
	sql.Register("pgxutc", utcer.Wrap(stdlib.GetDefaultDriver()))
}
