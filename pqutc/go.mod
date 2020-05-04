module github.com/hirofumi/utcer/pqutc

go 1.14

require (
	github.com/hirofumi/utcer v0.0.0-00010101000000-000000000000
	github.com/hirofumi/utcer/test v0.0.0-00010101000000-000000000000
	github.com/lib/pq v1.5.1
	github.com/stretchr/testify v1.5.1
)

replace (
	github.com/hirofumi/utcer => ../../utcer
	github.com/hirofumi/utcer/test => ../../utcer/test
)
