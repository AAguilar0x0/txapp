// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package psqldal

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID        pgtype.UUID `json:"id"`
	Email     string      `json:"email"`
	Password  string      `json:"password"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Role      string      `json:"role"`
}
