-- name: CreateOutboxMessage :exec

INSERT INTO outbox_messages (repo,owner) VALUES($1,$2);

-- name: SetSentStatusOutboxMessage :exec

UPDATE outbox_messages SET status = 'SENT' WHERE id = $1;

-- name: GetOutboxMessage :many

SELECT * FROM outbox_messages WHERE status = 'PENDING'
ORDER BY created_at
LIMIT 5
FOR UPDATE SKIP LOCKED;

-- name: GetRepositoryInfo :one

SELECT * FROM repositories WHERE fullname = LOWER($1);


-- name: CreateInboxMessage :exec

INSERT INTO inbox_messages (id,payload) VALUES($1,$2);


-- name: CreateOrUpdateRepoInfo :exec

INSERT INTO repositories (fullname,description,forks,stargazers,created_at,status) VALUES(LOWER($1),$2,$3,$4,$5,'READY') ON CONFLICT (fullname) DO UPDATE SET 
    updated_at = NOW(),
    description = EXCLUDED.description, 
    forks = EXCLUDED.forks, 
    stargazers = EXCLUDED.stargazers,
    status = EXCLUDED.status,
    created_at = EXCLUDED.created_at;


-- name: DeleteRepo :exec

DELETE FROM repositories WHERE fullname = LOWER($1);

-- name: CreateFetchingTask :exec

INSERT INTO repositories (fullname) VALUES(LOWER($1)); 

-- name: SetErrorStatusRepo :exec

UPDATE repositories SET status = 'ERROR' WHERE fullname = LOWER($1);