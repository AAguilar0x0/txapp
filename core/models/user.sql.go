package models

type GetUserForAuthRow struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
