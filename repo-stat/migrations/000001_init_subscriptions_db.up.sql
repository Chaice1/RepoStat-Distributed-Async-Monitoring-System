CREATE TABLE IF NOT EXISTS subscriptions(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repo_name VARCHAR(128) NOT NULL,
    owner_name VARCHAR(128) NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS outbox_messages(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repo VARCHAR(128) NOT NULL,
    owner VARCHAR(128) NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

