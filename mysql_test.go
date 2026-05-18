package data_mysql

import (
	"errors"
	"testing"

	mysql "github.com/go-sql-driver/mysql"
	"github.com/infrago/data"
)

func TestMySQLDialectClassifiesErrors(t *testing.T) {
	err := &mysql.MySQLError{Number: 1062, Message: "duplicate entry"}
	got := data.Error("insert", data.ErrInvalidUpdate, (mysqlDialect{}).ClassifyError(err))
	if !errors.Is(got, data.ErrDuplicate) {
		t.Fatalf("expected duplicate classification, got %v", got)
	}
	if !errors.Is(got, data.ErrConflict) {
		t.Fatalf("duplicate should be conflict-compatible")
	}
}
