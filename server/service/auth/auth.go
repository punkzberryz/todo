package auth

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	db "github.com/punkzberryz/todo/db/sqlc"
	"github.com/punkzberryz/todo/util"
)

var ErrEmailInUsed = fmt.Errorf("email has been used")

type Auth struct {
	Store db.Store
}

// Create user in database
type CreateUserParams struct {
	Email    string
	Username string
	Password string
}

func (a Auth) CreateUser(ctx context.Context, arg *CreateUserParams) (*db.User, error) {
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

func (a Auth) GetUserFromLogin(ctx context.Context, arg GetUserFromLogin) (*db.User, error) {
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
func (a Auth) GetUserFromToken(ctx context.Context, arg db.GetUserParams) (*db.User, error) {
	user, err := a.Store.GetUser(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
