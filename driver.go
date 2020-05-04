package utcer

import (
	"context"
	"database/sql/driver"
)

func wrapDriver(d driver.Driver) driver.Driver {
	if dc, ok := d.(driverContext); ok {
		return &driverContextUTC{driverContext: dc}
	}

	return &driverUTC{Driver: d}
}

type driverUTC struct {
	driver.Driver
}

func (d *driverUTC) Open(name string) (driver.Conn, error) {
	return openUTC(d.Driver, name)
}

type driverContext interface {
	driver.Driver
	driver.DriverContext
}

type driverContextUTC struct {
	driverContext
}

func (d *driverContextUTC) Open(name string) (driver.Conn, error) {
	return openUTC(d.driverContext, name)
}

func (d *driverContextUTC) OpenConnector(name string) (driver.Connector, error) {
	c, err := d.driverContext.OpenConnector(name)
	if err != nil {
		return nil, err
	}

	return &connectorUTC{Connector: c}, nil
}

type connectorUTC struct {
	driver.Connector
}

func (c *connectorUTC) Connect(ctx context.Context) (driver.Conn, error) {
	conn, err := c.Connector.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return wrapConn(conn), nil
}

func (c *connectorUTC) Driver() driver.Driver {
	return wrapDriver(c.Connector.Driver())
}

func openUTC(d driver.Driver, name string) (driver.Conn, error) {
	conn, err := d.Open(name)
	if err != nil {
		return nil, err
	}

	return wrapConn(conn), nil
}
