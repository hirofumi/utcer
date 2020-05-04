package utcer

import (
	"context"
	"database/sql/driver"
)

//go:generate mockgen -source=mock_test.go -destination mock_stmt_test.go -package utcer

// nolint: unused // interfaces for mockgen
type (
	Stmt interface {
		driver.Stmt
	}

	StmtWithChecker interface {
		driver.Stmt
		driver.NamedValueChecker
	}

	StmtExecContext interface {
		driver.Stmt
		driver.StmtExecContext
	}

	StmtExecContextWithChecker interface {
		driver.Stmt
		driver.StmtExecContext
		driver.NamedValueChecker
	}

	StmtQueryContext interface {
		driver.Stmt
		driver.StmtQueryContext
	}

	StmtQueryContextWithChecker interface {
		driver.Stmt
		driver.StmtQueryContext
		driver.NamedValueChecker
	}

	StmtContext interface {
		driver.Stmt
		driver.StmtExecContext
		driver.StmtQueryContext
	}

	StmtContextWithChecker interface {
		driver.Stmt
		driver.StmtExecContext
		driver.StmtQueryContext
		driver.NamedValueChecker
	}
)

type mockConnector struct {
	driver driver.Driver
}

func (m mockConnector) Connect(_ context.Context) (driver.Conn, error) {
	return m.driver.Open("")
}

func (m mockConnector) Driver() driver.Driver {
	return m.driver
}

type mockDriver struct {
	conn driver.Conn
}

func (d mockDriver) Open(_ string) (driver.Conn, error) {
	return d.conn, nil
}

type mockDriverContext struct {
	conn driver.Conn
}

func (d mockDriverContext) OpenConnector(_ string) (driver.Connector, error) {
	return mockConnector{driver: d}, nil
}

func (d mockDriverContext) Open(_ string) (driver.Conn, error) {
	return d.conn, nil
}

type mockConn struct {
	stmt driver.Stmt
}

func (c mockConn) Prepare(_ string) (driver.Stmt, error) {
	return c.stmt, nil
}

func (c mockConn) Close() error {
	return nil
}

func (c mockConn) Begin() (driver.Tx, error) {
	panic("implement me")
}

type mockConnPrepareContext struct {
	stmt driver.Stmt
}

func (c mockConnPrepareContext) PrepareContext(_ context.Context, _ string) (driver.Stmt, error) {
	return c.stmt, nil
}

func (c mockConnPrepareContext) Prepare(_ string) (driver.Stmt, error) {
	return c.stmt, nil
}

func (c mockConnPrepareContext) Close() error {
	return nil
}

func (c mockConnPrepareContext) Begin() (driver.Tx, error) {
	panic("implement me")
}
