package utcer

import "database/sql/driver"

func wrapStmt(s driver.Stmt) driver.Stmt {
	checker, ok := s.(driver.NamedValueChecker)
	if !ok {
		return s
	}

	checkerUTC := namedValueCheckerUTC{checker}

	switch s := s.(type) {
	case stmtContext:
		return &stmtContextUTC{s, checkerUTC}
	case stmtExecContext:
		return &stmtExecContextUTC{s, checkerUTC}
	case stmtQueryContext:
		return &stmtQueryContextUTC{s, checkerUTC}
	default:
		return &stmtUTC{s, checkerUTC}
	}
}

type stmtUTC struct {
	driver.Stmt
	namedValueCheckerUTC
}

type stmtExecContext interface {
	driver.Stmt
	driver.StmtExecContext
}

type stmtExecContextUTC struct {
	stmtExecContext
	namedValueCheckerUTC
}

type stmtQueryContext interface {
	driver.Stmt
	driver.StmtQueryContext
}

type stmtQueryContextUTC struct {
	stmtQueryContext
	namedValueCheckerUTC
}

type stmtContext interface {
	driver.Stmt
	driver.StmtExecContext
	driver.StmtQueryContext
}

type stmtContextUTC struct {
	stmtContext
	namedValueCheckerUTC
}
