-- name: CreateTransfer :execresult
INSERT INTO transfers (
    from_account_id,
    to_account_id,
    amount,
    created_at
) VALUES (
    ?, ?, ?, ?
);

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = ? LIMIT 1;

-- name: ListTransfer :many
SELECT * FROM transfers
WHERE from_account_id = ? OR
      to_account_id = ?
ORDER BY id LIMIT ?,?;