-- name: CreateUser :one
INSERT INTO users (
 email, first_name, last_name, password, role
) VALUES (
 $1, $2, $3, $4, $5
)
RETURNING *;

-- name: CountUser :one
SELECT count(*) FROM users;

