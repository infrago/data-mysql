package data_mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"

	mysql "github.com/go-sql-driver/mysql"
	. "github.com/infrago/base"
	"github.com/infrago/data"
)

type (
	mysqlDriver struct{}

	mysqlConnection struct {
		instance *data.Instance
		db       *sql.DB
		actives  int64
	}

	mysqlDialect struct{}
)

func (d *mysqlDriver) Connect(inst *data.Instance) (data.Connection, error) {
	return &mysqlConnection{instance: inst}, nil
}

func (c *mysqlConnection) Open() error {
	dsn := strings.TrimSpace(c.instance.Config.Url)
	if dsn == "" {
		if v, ok := c.instance.Setting["dsn"].(string); ok {
			dsn = v
		}
	}
	if dsn == "" {
		return fmt.Errorf("missing mysql dsn")
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return err
	}
	c.db = db
	return nil
}

func (c *mysqlConnection) Close() error {
	if c.db == nil {
		return nil
	}
	err := c.db.Close()
	c.db = nil
	return err
}

func (c *mysqlConnection) Health() data.Health {
	return data.Health{Workload: atomic.LoadInt64(&c.actives)}
}

func (c *mysqlConnection) DB() *sql.DB {
	return c.db
}

func (c *mysqlConnection) Dialect() data.Dialect {
	return mysqlDialect{}
}

func (mysqlDialect) Name() string { return "mysql" }
func (mysqlDialect) Quote(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "`", "")
	return "`" + s + "`"
}
func (mysqlDialect) Placeholder(_ int) string { return "?" }
func (mysqlDialect) SupportsILike() bool      { return false }
func (mysqlDialect) SupportsReturning() bool  { return false }
func (mysqlDialect) MaxParams() int           { return 65535 }
func (mysqlDialect) ClassifyError(err error) error {
	var myerr *mysql.MySQLError
	if !errors.As(err, &myerr) {
		return nil
	}
	switch myerr.Number {
	case 1062:
		return data.ErrDuplicate
	case 1451, 1452:
		return data.ErrForeignKey
	case 1205:
		return data.ErrTimeout
	case 1213:
		return data.ErrConflict
	case 1317:
		return data.ErrCanceled
	case 2006, 2013:
		return data.ErrDriver
	default:
		return nil
	}
}
func (mysqlDialect) BindValue(cfg Var, v any) (any, bool) {
	switch {
	case data.IsJSONVar(cfg):
		return data.BindJSONValue(v)
	case data.IsBinaryVar(cfg):
		return data.BindBinaryValue(v)
	case data.IsUUIDVar(cfg), data.IsDecimalVar(cfg):
		return data.BindTextValue(v)
	case data.IsTimeVar(cfg):
		return data.BindTimeValue(v)
	default:
		return nil, false
	}
}
func (mysqlDialect) DecodeValue(cfg Var, value any) (any, bool) {
	switch {
	case data.IsJSONVar(cfg):
		return data.DecodeJSONValue(value)
	case data.IsBinaryVar(cfg):
		return data.DecodeBinaryValue(value)
	case data.IsUUIDVar(cfg), data.IsDecimalVar(cfg):
		return data.DecodeTextValue(value)
	case data.IsTimeVar(cfg):
		return data.DecodeTimeValue(value)
	default:
		return nil, false
	}
}
