package utcer

import "database/sql/driver"

func Wrap(d driver.Driver) driver.Driver {
	return wrapDriver(d)
}
