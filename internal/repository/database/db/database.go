package db

import (
	"context"
	"fmt"

	"cashier/internal/pkg/errors"
	iDB "cashier/internal/repository/database"

	"database/sql/driver"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrNilTx = errors.New("tx is nil, begin first or use Transaction")

type database struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
	tx      *gorm.DB
}

func New(read, write *gorm.DB) iDB.IDatabase {
	return &database{
		readDB:  read,
		writeDB: write,
	}
}

func (db *database) WriteDB(ctx context.Context) *gorm.DB {
	return db.getWriteDB().WithContext(ctx)
}

func (db *database) ReadDB(ctx context.Context) *gorm.DB {
	return db.getReadDB().WithContext(ctx)
}

func (db *database) getWriteDB() *gorm.DB {
	if db.tx != nil {
		return db.tx
	}
	return db.writeDB
}

func (db *database) getReadDB() *gorm.DB {
	if db.tx != nil {
		return db.tx
	}
	return db.readDB
}

func (db *database) Begin(ctx context.Context) iDB.IDatabase {
	tx := db.writeDB.WithContext(ctx).Begin()
	return &database{
		tx:      tx,
		readDB:  db.readDB,
		writeDB: db.writeDB,
	}
}

func (db *database) Commit() error {
	if db.tx == nil {
		return ErrNilTx
	}

	return db.tx.Commit().Error
}

func (db *database) Rollback() error {
	if db.tx == nil {
		return ErrNilTx
	}

	return db.tx.Rollback().Error
}

func (db *database) Transaction(ctx context.Context, f func(context.Context, iDB.IDatabase) error) (txErr error) {
	txRepo := db.Begin(ctx)

	defer func() {
		r := recover()
		if r != nil {
			txErr = errors.Wrap(errors.ErrInternalServerError, fmt.Sprint(r))
		}
		if txErr != nil {
			_ = txRepo.(*database).Rollback()
		} else {
			_ = txRepo.(*database).Commit()
		}

	}()

	txErr = f(ctx, txRepo)
	if txErr != nil {
		return txErr
	}

	return nil
}

func notFoundOrInternalError(err error) error {
	if err == nil {
		return nil
	}
	if err == gorm.ErrRecordNotFound {
		return errors.ErrResourceNotFound
	}
	return errors.ErrInternalError
}

func duplicateOrInternalError(err error) error {
	if err == nil {
		return nil
	}
	switch e := err.(type) {
	case *mysql.MySQLError:
		if e.Number == 1062 {
			return errors.ErrResourceAlreadyExists
		}
	}
	return errors.ErrInternalError
}

// gormExpr 包裝 clause.Expr 才能讓 gorm 用來更新欄位
type gormExpr struct {
	clause.Expr
}

// GormValue ...
func (expr gormExpr) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return expr.Expr
}

// Value ...
func (expr gormExpr) Value() (driver.Value, error) {
	return expr.Expr, nil
}
