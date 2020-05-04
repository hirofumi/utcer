# UTCer

UTCer is a wrapper for database/sql/driver.Driver that sets time.Time location to UTC.

## Usage

### pgx

```go
import (
	"database/sql"

	_ "github.com/hirofumi/utcer/pgxutc"
)

func open(dsn string) (*sql.DB, error){
	return sql.Open("pgxutc", dsn)
}
```

### pq

```go
import (
	"database/sql"

	_ "github.com/hirofumi/utcer/pqutc"
)

func open(dsn string) (*sql.DB, error){
	return sql.Open("pqutc", dsn)
}
```
