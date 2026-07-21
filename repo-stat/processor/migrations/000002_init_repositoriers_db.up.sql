CREATE TABLE IF NOT EXISTS repositories(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    fullname VARCHAR(128) NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    forks INT NOT NULL DEFAULT 0,
    stargazers INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    status VARCHAR(16) NOT NULL DEFAULT 'FETCHING',
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS outbox_messages(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repo VARCHAR(128) NOT NULL,
    owner VARCHAR(128) NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS inbox_messages(
    id UUID PRIMARY KEY,
    Payload JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_outbox_pending ON outbox_messages(status) WHERE status = 'PENDING';
CREATE INDEX IF NOT EXISTS idx_repositories_fullname ON repositories(LOWER(fullname));
