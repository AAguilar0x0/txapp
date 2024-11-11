-- name: CreateUser :one
INSERT INTO users (
 id, email, first_name, last_name, password, role
) VALUES (
 $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetUserForAuth :one
SELECT users.id, users.password, users.role FROM users WHERE email = $1;
