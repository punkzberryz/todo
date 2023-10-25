package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/punkzberryz/todo/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomPasswordResetSession(t *testing.T, email string) PasswordResetSession {
	arg := CreatePasswordResetSessionParams{
		Email:     email,
		Otp:       util.RandomString(6),
		ExpiresAt: time.Now().Add(time.Minute),
	}

	session, err := testQueries.CreatePasswordResetSession(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, session)

	require.Equal(t, email, session.Email)
	require.Equal(t, arg.Otp, session.Otp)
	require.WithinDuration(t, arg.ExpiresAt, session.ExpiresAt, time.Second)
	require.NotEmpty(t, session.CreatedAt)
	return session
}

func TestCreatePasswordResetSession(t *testing.T) {
	user := CreateRandomUser(t)
	CreateRandomPasswordResetSession(t, user.Email)
}

func TestGetPasswordResetSession(t *testing.T) {
	user := CreateRandomUser(t)
	session := CreateRandomPasswordResetSession(t, user.Email)
	session2, err := testQueries.GetPasswordResetSession(context.Background(), user.Email)
	require.NoError(t, err)
	require.NotEmpty(t, session2)

	require.Equal(t, session.Email, session2.Email)
	require.Equal(t, session.Otp, session2.Otp)
	require.Equal(t, session.CreatedAt, session2.CreatedAt)
	require.Equal(t, session.ExpiresAt, session2.ExpiresAt)
}
func TestUpdatePasswordResetSession(t *testing.T) {
	user := CreateRandomUser(t)
	session := CreateRandomPasswordResetSession(t, user.Email)

	arg := UpdatePasswordResetSessionParams{
		Email:     user.Email,
		ExpiresAt: time.Now().Add(time.Hour),
		Otp:       util.RandomString(10),
	}
	session2, err := testQueries.UpdatePasswordResetSession(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, session2)
	require.Equal(t, session.Email, session2.Email)
	require.Equal(t, arg.Otp, session2.Otp)
	require.Equal(t, session.CreatedAt, session2.CreatedAt)
	require.WithinDuration(t, arg.ExpiresAt, session2.ExpiresAt, time.Second)
}
func TestUpdatePasswordResetSessionWithoutExistingSession(t *testing.T) {
	user := CreateRandomUser(t)
	arg := UpdatePasswordResetSessionParams{
		Email:     user.Email,
		Otp:       util.RandomString(6),
		ExpiresAt: time.Now().Add(time.Minute),
	}

	session, err := testQueries.UpdatePasswordResetSession(context.Background(), arg)
	require.NotEmpty(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, session)
}

func TestDeletePasswordResetSession(t *testing.T) {
	user := CreateRandomUser(t)
	CreateRandomPasswordResetSession(t, user.Email)
	err := testQueries.DeletePasswordResetSession(context.Background(), user.Email)
	require.NoError(t, err)
	session, err := testQueries.GetPasswordResetSession(context.Background(), user.Email)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, session)
}
func TestCreatePasswordResetSessionWithoutUserEmail(t *testing.T) {
	email := "some email that doesn't exist for sure..."
	arg := CreatePasswordResetSessionParams{
		Email:     email,
		Otp:       util.RandomString(6),
		ExpiresAt: time.Now().Add(time.Minute),
	}

	session, err := testQueries.CreatePasswordResetSession(context.Background(), arg)
	pqErr, ok := err.(*pq.Error)
	require.True(t, ok)
	require.Equal(t, pqErr.Code.Name(), "foreign_key_violation")
	require.Empty(t, session)
}
