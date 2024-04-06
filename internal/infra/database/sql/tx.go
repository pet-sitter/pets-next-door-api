package sql

import (
	"database/sql"
	"errors"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
)

type SqlTx struct {
	*sql.Tx
}

func (sct *SqlTx) EndTx(f func() *pnd.AppError) *pnd.AppError {
	var err *pnd.AppError
	tx := sct.Tx

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err := f(); err != nil {
			tx.Rollback()
		} else if err := tx.Commit(); err != nil {
			tx.Rollback()
		}
	}()

	err = f()
	return err
}

func (sct *SqlTx) Rollback() *pnd.AppError {
	if err := sct.Tx.Rollback(); err != nil {
		if errors.Is(err, sql.ErrTxDone) {
			return nil
		}
		return pnd.FromPostgresError(err)
	}

	return nil
}

func (sct *SqlTx) Commit() *pnd.AppError {
	if err := sct.Tx.Commit(); err != nil {
		return pnd.FromPostgresError(err)
	}

	return nil
}
