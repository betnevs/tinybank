-- name: CreateEntry :execresult
INSERT INTO entries (
    account_id,
    amount,
    created_at
) VALUES (
    ?, ?, ?
);

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = ? LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
WHERE account_id = ?
ORDER BY id LIMIT ?, ?;