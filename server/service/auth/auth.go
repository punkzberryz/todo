package auth

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	db "github.com/punkzberryz/todo/db/sqlc"
	"github.com/punkzberryz/todo/util"
)

var (
	ErrEmailInUsed   = fmt.Errorf("email has been used")
	ErrEmailNotFound = fmt.Errorf("email not found")
	ErrOtpNotMatched = fmt.Errorf("otp not matched")
	ErrOtpTimeOut    = fmt.Errorf("otp timeout")
)

type Auth struct {
	Store db.Store
}

// Create user in database
type CreateUserParams struct {
	Email    string
	Username string
	Password string
}

func (a *Auth) CreateUser(ctx context.Context, arg *CreateUserParams) (*db.User, error) {
	hashedPassword, err := util.HashPassword(arg.Password)
	if err != nil {
		return nil, err
	}

	user, err := a.Store.CreateUser(ctx, db.CreateUserParams{
		HashedPassword: hashedPassword,
		Username:       arg.Username,
		Email:          arg.Email,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, ErrEmailInUsed
			}
		}
		return nil, err
	}
	return &user, nil
}

// Get user from database
type GetUserFromLogin struct {
	db.GetUserParams
	Password string
}

func (a *Auth) GetUserFromLogin(ctx context.Context, arg GetUserFromLogin) (*db.User, error) {
	user, err := a.Store.GetUser(ctx, arg.GetUserParams)
	if err != nil {
		return nil, err
	}
	err = util.ComparePassword(arg.Password, user.HashedPassword)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Get User after JWT verified
func (a *Auth) GetUserFromToken(ctx context.Context, arg db.GetUserParams) (*db.User, error) {
	user, err := a.Store.GetUser(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Get ResetPassword OTP
func (a *Auth) GetResetPasswordOtp(ctx context.Context, email string) (*db.PasswordResetSession, error) {
	otp := fmt.Sprint(util.RandomInt(100000, 999999))
	session, err := a.Store.UpdatePasswordResetSession(ctx,
		db.UpdatePasswordResetSessionParams{
			Email:     email,
			Otp:       otp,
			ExpiresAt: time.Now().Add(time.Second * 61), //1 minute + 1 second
		})
	if err != nil && err != sql.ErrNoRows {
		log.Print("error here..", err)
		return nil, err
	}
	if err == nil {
		//session update succesfully
		return &session, nil
	}
	//session not exist, let's create one
	session, err = a.Store.CreatePasswordResetSession(ctx,
		db.CreatePasswordResetSessionParams{
			Email:     email,
			Otp:       otp,
			ExpiresAt: time.Now().Add(time.Second * 61), //1 minute + 1 second
		})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation":
				//fail because email in user table doesn't exist
				return nil, ErrEmailNotFound
			}
		}
		return nil, err
	}
	return &session, nil
}

// update user
type UpdateUserPasswordParam struct {
	Email       string
	Otp         string
	NewPassword string
}

func (a *Auth) UpdateUserPassword(ctx context.Context, arg *UpdateUserPasswordParam) error {
	//check if otp  is correct
	session, err := a.Store.GetPasswordResetSession(ctx, arg.Email)
	if err != nil {
		return err
	}

	if session.Otp != arg.Otp {
		return ErrOtpNotMatched
	}
	if session.ExpiresAt.Before(time.Now()) {
		return ErrOtpTimeOut
	}
	hashedPassword, err := util.HashPassword(arg.NewPassword)
	if err != nil {
		return err
	}

	_, err = a.Store.UpdateUser(ctx, db.UpdateUserParams{
		HashedPassword:    hashedPassword,
		PasswordChangedAt: time.Now(),
		Email:             session.Email,
		NewEmail:          session.Email,
	})
	//TODO: delete session after update sucess
	return err
}
