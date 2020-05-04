package test

import (
	"database/sql"
	"time"

	"github.com/stretchr/testify/suite"
)

type PostgresSuite struct {
	suite.Suite
	DriverName string
}

func (s *PostgresSuite) TestSelect() {
	db, err := sql.Open(s.DriverName, "host=localhost port=5432 user=postgres password=utcer sslmode=disable")
	s.Require().NoError(err)

	now := time.Now().Round(time.Microsecond).In(time.FixedZone("+09:00", int(9*time.Hour/time.Second)))
	row := db.QueryRow("SELECT CAST($1 AS TIMESTAMP)", now)

	var got time.Time
	s.NoError(row.Scan(&got))
	s.WithinDuration(now, got, 0)
}
