package models

type UserGetForAuthRow struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
