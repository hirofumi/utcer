# UTCer

UTCer is a wrapper for [`database/sql/driver.Driver`](https://pkg.go.dev/database/sql/driver@go1.14.2?tab=doc#Driver)
that sets [`time.Time`](https://pkg.go.dev/time@go1.14.2?tab=doc#Time) location to UTC.

## Motivation

[`.UTC()`](https://pkg.go.dev/time@go1.14.2?tab=doc#Time.UTC) appears again and again
in projects that store date and time as local `timestamp` in UTC.
UTCer mitigates this pain in exchange for inconvenience of `timestamp with time zone`.

## How It Works

UTCer implements and uses a [`driver.NamedValueChecker`](https://pkg.go.dev/database/sql/driver@go1.14.2?tab=doc#NamedValueChecker)
which replaces `Value` of [`driver.NamedValue`](https://pkg.go.dev/database/sql/driver@go1.14.2?tab=doc#NamedValue)
with `Value.UTC()` when its type is `time.Time`.

## Usage

### pgx

```go
import (
	"database/sql"

	_ "github.com/hirofumi/utcer/utcpgx"
)

func open(dsn string) (*sql.DB, error){
	return sql.Open("utcpgx", dsn)
}
```

### pq

```go
import (
	"database/sql"

	_ "github.com/hirofumi/utcer/utcpq"
)

func open(dsn string) (*sql.DB, error){
	return sql.Open("utcpq", dsn)
}
```
