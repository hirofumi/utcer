package test

import (
	"database/sql"
	"time"

	"github.com/stretchr/testify/suite"
)

type PostgresSuite struct {
	suite.Suite
	DriverName string
	db         *sql.DB
}

func (s *PostgresSuite) SetupSuite() {
	var err error

	s.db, err = sql.Open(s.DriverName, "host=localhost port=5432 user=postgres password=utcer sslmode=disable")
	s.Require().NoError(err)
}

func (s *PostgresSuite) TearDownSuite() {
	s.NoError(s.db.Close())
}

func (s *PostgresSuite) TestSelect() {
	now := time.Now().Round(time.Microsecond).In(time.FixedZone("+09:00", int(9*time.Hour/time.Second)))
	row := s.db.QueryRow("SELECT CAST($1 AS TIMESTAMP)", now)

	var got time.Time
	s.NoError(row.Scan(&got))
	s.WithinDuration(now, got, 0)
}
