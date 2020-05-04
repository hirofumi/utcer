package utcer

import (
	"context"
	"database/sql/driver"
)

func wrapConn(conn driver.Conn) driver.Conn {
	checker, ok := conn.(driver.NamedValueChecker)
	if !ok {
		checker = namedValueCheckerNOP{}
	}

	checkerUTC := namedValueCheckerUTC{checker}

	if conn, ok := conn.(connPrepareContext); ok {
		return &connPrepareContextUTC{conn, checkerUTC}
	}

	return &connUTC{conn, checkerUTC}
}

type connUTC struct {
	driver.Conn
	namedValueCheckerUTC
}

func (c *connUTC) Prepare(query string) (driver.Stmt, error) {
	return prepareUTC(c.Conn, query)
}

type connPrepareContext interface {
	driver.Conn
	driver.ConnPrepareContext
}

type connPrepareContextUTC struct {
	connPrepareContext
	namedValueCheckerUTC
}

func (c *connPrepareContextUTC) Prepare(query string) (driver.Stmt, error) {
	return prepareUTC(c.connPrepareContext, query)
}

func (c *connPrepareContextUTC) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	s, err := c.connPrepareContext.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return wrapStmt(s), nil
}

func prepareUTC(conn driver.Conn, query string) (driver.Stmt, error) {
	s, err := conn.Prepare(query)
	if err != nil {
		return nil, err
	}

	return wrapStmt(s), nil
}
