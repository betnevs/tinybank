-- name: CreateUser :execresult
INSERT INTO users (
    username,
    hashed_password,
    full_name,
    email
) VALUES (
    ?, ?, ?, ?
);

-- name: GetUser :one
SELECT * FROM users
WHERE username = ? LIMIT 1;