-- name: CreateSubscription :exec

INSERT INTO subscriptions (repo_name,owner_name) VALUES($1,$2);

-- name: DeleteSubscription :exec

DELETE FROM subscriptions WHERE repo_name = $1 AND owner_name = $2;

-- name: GetSubscriptions :many

SELECT * FROM subscriptions;

-- name: GetSubscription :one

SELECT * FROM subscriptions WHERE repo_name = $1 AND owner_name = $2;

-- name: GetOutboxMessage :one

SELECT * FROM outbox_messages WHERE status = 'PENDING'
ORDER BY created_at
LIMIT 1
FOR UPDATE SKIP LOCKED;

-- name: SetSentStatusOutboxMessage :exec

UPDATE outbox_messages SET status = 'SENT' WHERE id = $1;

-- name: CreateOutboxMessage :exec

INSERT INTO outbox_messages (repo,owner) VALUES($1,$2);
