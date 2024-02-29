package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func (sdt *SqlxDBTx) WithinTransaction(ctx context.Context, txFunc func(ctx context.Context) error) error {
	return sdt.WithinTransactionOptions(ctx, txFunc, nil)
}

func (sdt *SqlxDBTx) WithinTransactionOptions(ctx context.Context, txFunc func(ctx context.Context) error, txOptions *sql.TxOptions) error {
	var err error
	var tx *sqlx.Tx

	if txOptions != nil {
		tx, err = sdt.DB.BeginTxx(ctx, txOptions)
	} else {
		tx, err = sdt.DB.Beginx()
	}

	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	sct := &SqlxConnTx{tx}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = txFunc(InjectTx(ctx, sct))
	return err
}

func (sct *SqlxConnTx) WithinTransaction(ctx context.Context, txFunc func(ctx context.Context) error) error {
	return sct.WithinTransactionOptions(ctx, txFunc, nil)
}

func (sct *SqlxConnTx) WithinTransactionOptions(ctx context.Context, txFunc func(ctx context.Context) error, txOptions *sql.TxOptions) error {
	var err error
	tx := sct.DB
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = txFunc(InjectTx(ctx, sct))
	return err
}
