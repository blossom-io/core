package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"core/pkg/postgres"
)

type repository struct {
	DB *postgres.Postgres
}

type Transactor interface {
	InTX(ctx context.Context, txFunc []func(ctx context.Context) error) error
}

type Repository interface {
	Transactor
	User
	Subchater
}

var _ Repository = (*repository)(nil)

type txKey struct{}

func injectTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func extractTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}

	return nil
}

// New -.
func New(pg *postgres.Postgres) Repository {
	return &repository{DB: pg}
}

func (r *repository) prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.PrepareContext(ctx, query)
	}

	return r.DB.DB.PrepareContext(ctx, query)
}

func (r *repository) InTX(ctx context.Context, txFunc []func(ctx context.Context) error) error {
	// begin transaction
	tx, err := r.DB.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// run callback
	// err = txFunc(injectTx(ctx, tx))
	// if err != nil {
	// 	// if error, rollback
	// 	if errRollback := tx.Rollback(); errRollback != nil {
	// 		log.Printf("rollback transaction: %v", errRollback)
	// 	}
	// 	return err
	// }

	for _, f := range txFunc {
		err = f(injectTx(ctx, tx))
		if err != nil {
			// if error, rollback
			if errRollback := tx.Rollback(); errRollback != nil {
				log.Printf("rollback transaction: %v", errRollback)
			}
			return err
		}
	}

	// if no error, commit
	if errCommit := tx.Commit(); errCommit != nil {
		log.Printf("commit transaction: %v", errCommit)
	}
	return nil
}
