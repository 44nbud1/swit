package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
	"time"
)

func (r *Repository) CreateUser(ctx context.Context, input UserRequest) (UserResponse, error) {

	// initialize variable
	var err error
	var tag pgconn.CommandTag

	if err != nil {
		log.Printf("failed to generate uuid: %v", err)

		return UserResponse{}, err
	}

	if tag, err = r.Db.Exec(ctx, "INSERT INTO users (id, full_name, phone_number, created_at, updated_at) VALUES ($1,$2,$3,$4,$5)",
		input.ID,
		input.FullName,
		input.PhoneNumber,
		time.Now(),
		time.Now(),
	); err != nil {
		log.Printf("error when store data, err: %v", err)

		return UserResponse{}, err
	}

	defaultTag := 1
	if tag.RowsAffected() != int64(defaultTag) {
		log.Printf("failed to insert topup transaction, tag: %v", tag.RowsAffected())
	}

	return UserResponse{
		ID: input.ID,
	}, err
}

func (r *Repository) CreatePassword(ctx context.Context, input PasswordRequest) (PasswordResponse, error) {

	// initialize variable
	var err error
	var tag pgconn.CommandTag

	// uuid
	newUUID, err := uuid.NewUUID()
	if err != nil {
		log.Printf("failed to generate uuid: %v", err)

		return PasswordResponse{}, err
	}

	if tag, err = r.Db.Exec(ctx, "INSERT INTO password (id, password, user_id, created_at, updated_at) VALUES ($1,$2,$3,$4,$5)",
		newUUID.String(),
		input.Password,
		input.UserID,
		time.Now(),
		time.Now(),
	); err != nil {
		log.Printf("error when store data, err: %v", err)

		return PasswordResponse{}, err
	}

	defaultTag := 1
	if tag.RowsAffected() != int64(defaultTag) {
		log.Printf("failed to insert topup transaction, tag: %v", tag.RowsAffected())
	}

	return PasswordResponse{
		ID: newUUID.String(),
	}, err
}

func (r *Repository) InquiryUserByPhoneNumber(ctx context.Context, input UserRequest) (UserResponse, error) {

	var err error
	var userResponse UserResponse

	if err = r.Db.QueryRow(ctx,
		"SELECT id FROM users WHERE phone_number = $1", input.PhoneNumber).Scan(&userResponse.ID); err != nil {

		if err.Error() == "no rows in result set" {
			return userResponse, nil
		}
	}

	return userResponse, err
}

func (r *Repository) InquiryPasswordByUserID(ctx context.Context, input PasswordRequest) (PasswordResponse, error) {

	var err error
	var pwd PasswordResponse

	if err = r.Db.QueryRow(ctx,
		"SELECT password FROM password WHERE user_id = $1", input.UserID).Scan(&pwd.Password); err != nil {

		if err.Error() == "no rows in result set" {
			return pwd, nil
		}

		log.Printf("error when query password by id and pwd: %v", err)

	}

	return pwd, err
}

func (r *Repository) InquiryByPhoneNumber(ctx context.Context, phoneNumber string) (UserResponse, error) {

	var err error
	var userResponse UserResponse

	if err = r.Db.QueryRow(ctx,
		"SELECT phone_number, full_name FROM users WHERE phone_number= $1", phoneNumber).Scan(&userResponse.PhoneNumber, &userResponse.FullName); err != nil {
		log.Printf("error when query by id, err: %v", err)

		if errors.Is(err, sql.ErrNoRows) {
			return UserResponse{}, nil
		}

		return UserResponse{}, err
	}

	return userResponse, err
}
func (r *Repository) InquiryByUserID(ctx context.Context, id string) (UserResponse, error) {

	var err error
	var userResponse UserResponse

	if err = r.Db.QueryRow(ctx,
		"SELECT phone_number, full_name FROM users WHERE id= $1", id).Scan(&userResponse.PhoneNumber, &userResponse.FullName); err != nil {
		log.Printf("error when query by id, err: %v", err)

		return UserResponse{}, fmt.Errorf("data_not_found")
	}

	return userResponse, err
}

func (r *Repository) UpdateUser(ctx context.Context, input UserRequest) (UserResponse, error) {

	var err error
	var tag pgconn.CommandTag

	if tag, err = r.Db.Exec(ctx,
		`update users set phone_number=$1, full_name=$2 where id=$3`, input.PhoneNumber, input.FullName, input.ID); err != nil {
		log.Printf("error when query by id, err: %v", err)

		return UserResponse{}, fmt.Errorf("data_not_found")
	}

	defaultTag := 1
	if tag.RowsAffected() != int64(defaultTag) {
		log.Printf("failed to insert topup transaction, tag: %v", tag.RowsAffected())
	}

	return UserResponse{
		PhoneNumber: input.PhoneNumber,
		FullName:    input.FullName,
	}, err
}

func (r *Repository) CreateLogin(ctx context.Context, input LoginRequest) error {
	var err error
	var tag pgconn.CommandTag

	if tag, err = r.Db.Exec(ctx,
		"INSERT INTO login (id, user_id, created_at) VALUES ($1,$2,$3)", input.ID, input.UserID, time.Now()); err != nil {
		log.Printf("error when save data, err: %v", err)

		return err
	}

	defaultTag := 1
	if tag.RowsAffected() != int64(defaultTag) {
		log.Printf("failed to insert topup transaction, tag: %v", tag.RowsAffected())
	}

	return err
}
