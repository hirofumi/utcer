package utcer

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// nolint: funlen
func TestWrap(t *testing.T) {
	var (
		zoneNonUTC = time.FixedZone("+09:00", int(9*time.Hour/time.Second))
		timeUTC    = time.Date(2012, 3, 14, 0, 0, 0, 0, time.UTC)
		timeNonUTC = timeUTC.In(zoneNonUTC)
		valNonUTC  = timeValuer{time: timeNonUTC}
	)

	ctx := context.Background()

	var (
		e = func(v driver.Value) func(db *sql.DB) error {
			return func(db *sql.DB) error {
				_, err := db.Exec(`INSERT INTO table VALUES (?)`, v)

				return err
			}
		}
		ec = func(v driver.Value) func(db *sql.DB) error {
			return func(db *sql.DB) error {
				_, err := db.ExecContext(ctx, `INSERT INTO table VALUES (?)`, v)

				return err
			}
		}
		q = func(v driver.Value) func(db *sql.DB) error {
			return func(db *sql.DB) error {
				rows, err := db.Query(`SELECT ?`, v)
				if err != nil {
					return err
				}

				return rows.Err()
			}
		}
		qc = func(v driver.Value) func(db *sql.DB) error {
			return func(db *sql.DB) error {
				rows, err := db.QueryContext(ctx, `SELECT ?`, v)
				if err != nil {
					return err
				}

				return rows.Err()
			}
		}
		pcE = func(v driver.Value) func(db *sql.DB) error {
			return func(db *sql.DB) error {
				s, err := db.PrepareContext(ctx, `INSERT INTO table VALUES (?)`)
				if err != nil {
					return err
				}

				_, err = s.Exec(v)
				if err != nil {
					return err
				}

				return s.Close()
			}
		}
		pcEC = func(v driver.Value) func(db *sql.DB) error {
			return func(db *sql.DB) error {
				s, err := db.PrepareContext(ctx, `INSERT INTO table VALUES (?)`)
				if err != nil {
					return err
				}

				_, err = s.ExecContext(ctx, v)
				if err != nil {
					return err
				}

				return s.Close()
			}
		}
		pcQ = func(v driver.Value) func(db *sql.DB) error {
			return func(db *sql.DB) error {
				s, err := db.PrepareContext(ctx, `SELECT ?`)
				if err != nil {
					return err
				}

				rows, err := s.Query(v)
				if err != nil {
					return err
				}

				return rows.Err()
			}
		}
		pcQC = func(v driver.Value) func(db *sql.DB) error {
			return func(db *sql.DB) error {
				s, err := db.PrepareContext(ctx, `SELECT ?`)
				if err != nil {
					return err
				}

				rows, err := s.QueryContext(ctx, v)
				if err != nil {
					return err
				}

				return rows.Err()
			}
		}
	)

	var (
		sE = func(v driver.Value) func(*gomock.Controller) driver.Stmt {
			return func(ctrl *gomock.Controller) driver.Stmt {
				mock := NewMockStmt(ctrl)
				mock.EXPECT().NumInput().Return(1)
				mock.EXPECT().Exec([]driver.Value{v}).Return(nil, nil)
				mock.EXPECT().Close().Return(nil)
				return mock
			}
		}
		sQ = func(v driver.Value) func(*gomock.Controller) driver.Stmt {
			return func(ctrl *gomock.Controller) driver.Stmt {
				mock := NewMockStmt(ctrl)
				mock.EXPECT().NumInput().Return(1)
				mock.EXPECT().Query([]driver.Value{v}).Return(nil, nil)
				return mock
			}
		}
		swcE = func(v driver.Value) func(*gomock.Controller) driver.Stmt {
			return func(ctrl *gomock.Controller) driver.Stmt {
				mock := NewMockStmtWithChecker(ctrl)
				mock.EXPECT().NumInput().Return(1)
				mock.EXPECT().CheckNamedValue(&driver.NamedValue{Ordinal: 1, Value: v}).Return(nil)
				mock.EXPECT().Exec([]driver.Value{v}).Return(nil, nil)
				mock.EXPECT().Close().Return(nil)
				return mock
			}
		}
		swcQ = func(v driver.Value) func(*gomock.Controller) driver.Stmt {
			return func(ctrl *gomock.Controller) driver.Stmt {
				mock := NewMockStmtWithChecker(ctrl)
				mock.EXPECT().NumInput().Return(1)
				mock.EXPECT().CheckNamedValue(&driver.NamedValue{Ordinal: 1, Value: v}).Return(nil)
				mock.EXPECT().Query([]driver.Value{v}).Return(nil, nil)
				return mock
			}
		}
		secEC = func(v driver.Value) func(*gomock.Controller) driver.Stmt {
			return func(ctrl *gomock.Controller) driver.Stmt {
				mock := NewMockStmtExecContext(ctrl)
				mock.EXPECT().NumInput().Return(1)
				mock.EXPECT().ExecContext(ctx, []driver.NamedValue{{Ordinal: 1, Value: v}}).Return(nil, nil)
				mock.EXPECT().Close().Return(nil)
				return mock
			}
		}
		secQ = func(v driver.Value) func(*gomock.Controller) driver.Stmt {
			return func(ctrl *gomock.Controller) driver.Stmt {
				mock := NewMockStmtExecContext(ctrl)
				mock.EXPECT().NumInput().Return(1)
				mock.EXPECT().Query([]driver.Value{v}).Return(nil, nil)
				return mock
			}
		}
		secwcEC = func(v driver.Value) func(*gomock.Controller) driver.Stmt {
			return func(ctrl *gomock.Controller) driver.Stmt {
				mock := NewMockStmtExecContextWithChecker(ctrl)
				mock.EXPECT().NumInput().Return(1)
				mock.EXPECT().CheckNamedValue(&driver.NamedValue{Ordinal: 1, Value: v}).Return(nil)
				mock.EXPECT().ExecContext(ctx, []driver.NamedValue{{Ordinal: 1, Value: v}}).Return(nil, nil)
				mock.EXPECT().Close().Return(nil)
				return mock
			}
		}
		secwcQ = func(v driver.Value) func(*gomock.Controller) driver.Stmt {
			return func(ctrl *gomock.Controller) driver.Stmt {
				mock := NewMockStmtExecContextWithChecker(ctrl)
				mock.EXPECT().NumInput().Return(1)
				mock.EXPECT().CheckNamedValue(&driver.NamedValue{Ordinal: 1, Value: v}).Return(nil)
				mock.EXPECT().Query([]driver.Value{v}).Return(nil, nil)
				return mock
			}
		}
		sqcwcE = func(v driver.Value) func(*gomock.Controller) driver.Stmt {
			return func(ctrl *gomock.Controller) driver.Stmt {
				mock := NewMockStmtQueryContextWithChecker(ctrl)
				mock.EXPECT().NumInput().Return(1)
				mock.EXPECT().CheckNamedValue(&driver.NamedValue{Ordinal: 1, Value: v}).Return(nil)
				mock.EXPECT().Exec([]driver.Value{v}).Return(nil, nil)
				mock.EXPECT().Close().Return(nil)
				return mock
			}
		}
		sqcwcQC = func(v driver.Value) func(*gomock.Controller) driver.Stmt {
			return func(ctrl *gomock.Controller) driver.Stmt {
				mock := NewMockStmtQueryContextWithChecker(ctrl)
				mock.EXPECT().NumInput().Return(1)
				mock.EXPECT().CheckNamedValue(&driver.NamedValue{Ordinal: 1, Value: v}).Return(nil)
				mock.EXPECT().QueryContext(ctx, []driver.NamedValue{{Ordinal: 1, Value: v}}).Return(nil, nil)
				return mock
			}
		}
		scwcEC = func(v driver.Value) func(*gomock.Controller) driver.Stmt {
			return func(ctrl *gomock.Controller) driver.Stmt {
				mock := NewMockStmtContextWithChecker(ctrl)
				mock.EXPECT().NumInput().Return(1)
				mock.EXPECT().CheckNamedValue(&driver.NamedValue{Ordinal: 1, Value: v}).Return(nil)
				mock.EXPECT().ExecContext(ctx, []driver.NamedValue{{Ordinal: 1, Value: v}}).Return(nil, nil)
				mock.EXPECT().Close().Return(nil)
				return mock
			}
		}
		scwcQC = func(v driver.Value) func(*gomock.Controller) driver.Stmt {
			return func(ctrl *gomock.Controller) driver.Stmt {
				mock := NewMockStmtContextWithChecker(ctrl)
				mock.EXPECT().NumInput().Return(1)
				mock.EXPECT().CheckNamedValue(&driver.NamedValue{Ordinal: 1, Value: v}).Return(nil)
				mock.EXPECT().QueryContext(ctx, []driver.NamedValue{{Ordinal: 1, Value: v}}).Return(nil, nil)
				return mock
			}
		}
	)

	type test struct {
		name string
		mock func(ctrl *gomock.Controller) driver.Stmt
		run  func(db *sql.DB) error
	}

	run := func(name string, tt test, newMockDriver func(*gomock.Controller) driver.Driver) {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			assert.NoError(t, tt.run(sql.OpenDB(mockConnector{driver: Wrap(newMockDriver(ctrl))})))
		})
	}

	tests := []test{
		{name: "Stmt/Exec/Time", mock: sE(timeUTC), run: e(timeNonUTC)},
		{name: "Stmt/Exec/Valuer", mock: sE(timeNonUTC), run: e(valNonUTC)},
		{name: "Stmt/ExecContext/Time", mock: sE(timeUTC), run: ec(timeNonUTC)},
		{name: "Stmt/ExecContext/Valuer", mock: sE(timeNonUTC), run: e(valNonUTC)},
		{name: "Stmt/Query/Time", mock: sQ(timeUTC), run: q(timeNonUTC)},
		{name: "Stmt/Query/Valuer", mock: sQ(timeNonUTC), run: q(valNonUTC)},
		{name: "Stmt/QueryContext/Time", mock: sQ(timeUTC), run: qc(timeNonUTC)},
		{name: "Stmt/QueryContext/Valuer", mock: sQ(timeNonUTC), run: qc(valNonUTC)},
		{name: "StmtWithChecker/Exec/Time", mock: swcE(timeUTC), run: e(timeNonUTC)},
		{name: "StmtWithChecker/Exec/Valuer", mock: swcE(valNonUTC), run: e(valNonUTC)},
		{name: "StmtWithChecker/ExecContext/Time", mock: swcE(timeUTC), run: ec(timeNonUTC)},
		{name: "StmtWithChecker/ExecContext/Valuer", mock: swcE(valNonUTC), run: ec(valNonUTC)},
		{name: "StmtWithChecker/Query/Time", mock: swcQ(timeUTC), run: q(timeNonUTC)},
		{name: "StmtWithChecker/Query/Valuer", mock: swcQ(valNonUTC), run: q(valNonUTC)},
		{name: "StmtWithChecker/QueryContext/Time", mock: swcQ(timeUTC), run: qc(timeNonUTC)},
		{name: "StmtWithChecker/QueryContext/Valuer", mock: swcQ(valNonUTC), run: qc(valNonUTC)},
		{name: "StmtExecContext/Exec/Time", mock: secEC(timeUTC), run: e(timeNonUTC)},
		{name: "StmtExecContext/Exec/Valuer", mock: secEC(timeNonUTC), run: e(valNonUTC)},
		{name: "StmtExecContext/ExecContext/Time", mock: secEC(timeUTC), run: ec(timeNonUTC)},
		{name: "StmtExecContext/ExecContext/Valuer", mock: secEC(timeNonUTC), run: e(valNonUTC)},
		{name: "StmtExecContext/Query/Time", mock: secQ(timeUTC), run: q(timeNonUTC)},
		{name: "StmtExecContext/Query/Valuer", mock: secQ(timeNonUTC), run: q(valNonUTC)},
		{name: "StmtExecContext/QueryContext/Time", mock: secQ(timeUTC), run: qc(timeNonUTC)},
		{name: "StmtExecContext/QueryContext/Valuer", mock: secQ(timeNonUTC), run: qc(valNonUTC)},
		{name: "StmtExecContextWithChecker/Exec/Time", mock: secwcEC(timeUTC), run: e(timeNonUTC)},
		{name: "StmtExecContextWithChecker/Exec/Valuer", mock: secwcEC(valNonUTC), run: e(valNonUTC)},
		{name: "StmtExecContextWithChecker/ExecContext/Time", mock: secwcEC(timeUTC), run: ec(timeNonUTC)},
		{name: "StmtExecContextWithChecker/ExecContext/Valuer", mock: secwcEC(valNonUTC), run: ec(valNonUTC)},
		{name: "StmtExecContextWithChecker/Query/Time", mock: secwcQ(timeUTC), run: q(timeNonUTC)},
		{name: "StmtExecContextWithChecker/Query/Valuer", mock: secwcQ(valNonUTC), run: q(valNonUTC)},
		{name: "StmtExecContextWithChecker/QueryContext/Time", mock: secwcQ(timeUTC), run: qc(timeNonUTC)},
		{name: "StmtExecContextWithChecker/QueryContext/Valuer", mock: secwcQ(valNonUTC), run: qc(valNonUTC)},
		{name: "StmtQueryContextWithChecker/Exec/Time", mock: sqcwcE(timeUTC), run: e(timeNonUTC)},
		{name: "StmtQueryContextWithChecker/Exec/Valuer", mock: sqcwcE(valNonUTC), run: e(valNonUTC)},
		{name: "StmtQueryContextWithChecker/ExecContext/Time", mock: sqcwcE(timeUTC), run: ec(timeNonUTC)},
		{name: "StmtQueryContextWithChecker/ExecContext/Valuer", mock: sqcwcE(valNonUTC), run: ec(valNonUTC)},
		{name: "StmtQueryContextWithChecker/Query/Time", mock: sqcwcQC(timeUTC), run: q(timeNonUTC)},
		{name: "StmtQueryContextWithChecker/Query/Valuer", mock: sqcwcQC(valNonUTC), run: q(valNonUTC)},
		{name: "StmtQueryContextWithChecker/QueryContext/Time", mock: sqcwcQC(timeUTC), run: qc(timeNonUTC)},
		{name: "StmtQueryContextWithChecker/QueryContext/Valuer", mock: sqcwcQC(valNonUTC), run: qc(valNonUTC)},
		{name: "StmtContextWithChecker/Exec/Time", mock: scwcEC(timeUTC), run: e(timeNonUTC)},
		{name: "StmtContextWithChecker/Exec/Valuer", mock: scwcEC(valNonUTC), run: e(valNonUTC)},
		{name: "StmtContextWithChecker/ExecContext/Time", mock: scwcEC(timeUTC), run: ec(timeNonUTC)},
		{name: "StmtContextWithChecker/ExecContext/Valuer", mock: scwcEC(valNonUTC), run: ec(valNonUTC)},
		{name: "StmtContextWithChecker/Query/Time", mock: scwcQC(timeUTC), run: q(timeNonUTC)},
		{name: "StmtContextWithChecker/Query/Valuer", mock: scwcQC(valNonUTC), run: q(valNonUTC)},
		{name: "StmtContextWithChecker/QueryContext/Time", mock: scwcQC(timeUTC), run: qc(timeNonUTC)},
		{name: "StmtContextWithChecker/QueryContext/Valuer", mock: scwcQC(valNonUTC), run: qc(valNonUTC)},
		{name: "Stmt/PreparedContext/Exec/Time", mock: sE(timeUTC), run: pcE(timeNonUTC)},
		{name: "Stmt/PreparedContext/Exec/Valuer", mock: sE(timeNonUTC), run: pcE(valNonUTC)},
		{name: "Stmt/PreparedContext/ExecContext/Time", mock: sE(timeUTC), run: pcEC(timeNonUTC)},
		{name: "Stmt/PreparedContext/ExecContext/Valuer", mock: sE(timeNonUTC), run: pcE(valNonUTC)},
		{name: "Stmt/PreparedContext/Query/Time", mock: sQ(timeUTC), run: pcQ(timeNonUTC)},
		{name: "Stmt/PreparedContext/Query/Valuer", mock: sQ(timeNonUTC), run: pcQ(valNonUTC)},
		{name: "Stmt/PreparedContext/QueryContext/Time", mock: sQ(timeUTC), run: pcQC(timeNonUTC)},
		{name: "Stmt/PreparedContext/QueryContext/Valuer", mock: sQ(timeNonUTC), run: pcQC(valNonUTC)},
	}
	for i := range tests {
		tt := tests[i]

		run("Driver/Conn/"+tt.name, tt, func(ctrl *gomock.Controller) driver.Driver {
			return mockDriver{conn: mockConn{stmt: tt.mock(ctrl)}}
		})
		run("Driver/ConnPrepareContext/"+tt.name, tt, func(ctrl *gomock.Controller) driver.Driver {
			return mockDriver{conn: mockConnPrepareContext{stmt: tt.mock(ctrl)}}
		})
	}

	t.Run("DriverContext", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sql.Register("utcer-test", Wrap(mockDriverContext{conn: mockConn{stmt: sQ(timeUTC)(ctrl)}}))
		db, err := sql.Open("utcer-test", "")
		if err != nil {
			assert.NoError(t, err)
		}

		rows, err := db.Query(`SELECT ?`, timeNonUTC)
		if assert.NoError(t, err) {
			assert.NoError(t, rows.Err())
		}
	})
}

type timeValuer struct {
	time time.Time
}

func (v timeValuer) Value() (driver.Value, error) {
	return v.time, nil
}

func (v timeValuer) String() string {
	return fmt.Sprintf("timeValuer(%s)", v.time)
}
