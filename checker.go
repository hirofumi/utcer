package utcer

import (
	"database/sql/driver"
	"time"
)

type namedValueCheckerUTC struct {
	driver.NamedValueChecker
}

func (c namedValueCheckerUTC) CheckNamedValue(nv *driver.NamedValue) error {
	if t, ok := nv.Value.(time.Time); ok {
		nv.Value = t.UTC()
	}

	return c.NamedValueChecker.CheckNamedValue(nv)
}

type namedValueCheckerNOP struct{}

func (c namedValueCheckerNOP) CheckNamedValue(_ *driver.NamedValue) error {
	return driver.ErrSkip
}
