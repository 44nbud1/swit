// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	CreateUser(ctx context.Context, input UserRequest) (UserResponse, error)
	CreatePassword(ctx context.Context, input PasswordRequest) (PasswordResponse, error)
	InquiryUserByPhoneNumber(ctx context.Context, input UserRequest) (UserResponse, error)
	InquiryByUserID(ctx context.Context, id string) (UserResponse, error)
	UpdateUser(ctx context.Context, input UserRequest) (UserResponse, error)
	InquiryByPhoneNumber(ctx context.Context, phoneNumber string) (UserResponse, error)
	CreateLogin(ctx context.Context, input LoginRequest) error
	InquiryPasswordByUserID(ctx context.Context, input PasswordRequest) (PasswordResponse, error)
}

type Tx struct {
	connPool *pgxpool.Pool
	*Repository
}

func (t *Tx) RegistrationTx(ctx context.Context, u UserRequest, p PasswordRequest) (UserResponse, error) {
	var res UserResponse

	err := t.execTrx(ctx, func(repository *Repository) error {
		var err error

		if res, err = repository.CreateUser(ctx, u); err != nil {
			return err
		}

		if _, err = repository.CreatePassword(ctx, p); err != nil {
			return err
		}

		return err
	})

	return res, err
}

type UserTxRepo interface {
	UserRepository
	RegistrationTx(ctx context.Context, u UserRequest, p PasswordRequest) (UserResponse, error)
}

func NewTx(connPool *pgxpool.Pool) *Tx {
	return &Tx{
		connPool:   connPool,
		Repository: NewRepository(connPool),
	}
}

func (t *Tx) execTrx(ctx context.Context, fn func(*Repository) error) error {
	tx, err := t.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := NewRepository(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
