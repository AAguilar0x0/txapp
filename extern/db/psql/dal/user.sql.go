// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package dal

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
 id, email, first_name, last_name, password, role
) VALUES (
 $6, $1, $2, $3, $4, $5
)
RETURNING id, email, password, first_name, last_name, role, created_at, updated_at
`

func (q *Queries) CreateUser(ctx context.Context, email string, firstName string, lastName string, password string, role string, newID string) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		email,
		firstName,
		lastName,
		password,
		role,
		newID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.FirstName,
		&i.LastName,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserForAuth = `-- name: GetUserForAuth :one
SELECT users.id, users.password, users.role FROM users WHERE email = $1
`

type GetUserForAuthRow struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (q *Queries) GetUserForAuth(ctx context.Context, email string) (GetUserForAuthRow, error) {
	row := q.db.QueryRow(ctx, getUserForAuth, email)
	var i GetUserForAuthRow
	err := row.Scan(&i.ID, &i.Password, &i.Role)
	return i, err
}
