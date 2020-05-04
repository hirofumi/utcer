package pqutc

import (
	"testing"

	"github.com/hirofumi/utcer/test"
	"github.com/stretchr/testify/suite"
)

func TestPostgresSuite(t *testing.T) {
	suite.Run(t, &test.PostgresSuite{DriverName: "pqutc"})
}
