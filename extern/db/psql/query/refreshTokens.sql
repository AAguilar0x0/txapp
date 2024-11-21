-- name: RefreshTokenCreate :exec
INSERT INTO refresh_tokens (
  id, user_id, expires_at
) VALUES (
  $1, $2, $3
);

-- name: RefreshTokenGet :one
SELECT *
FROM refresh_tokenS rt
WHERE id = $1
AND user_id = $2;

-- name: RefreshTokenDelete :execrows
DELETE FROM refresh_tokens
WHERE id = $1;

-- name: RefreshTokenDeleteFromUser :execrows
DELETE FROM refresh_tokens
WHERE user_id = $1;
