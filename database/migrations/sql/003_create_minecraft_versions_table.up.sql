-- Create minecraft_versions table
CREATE TABLE IF NOT EXISTS minecraft_versions (
    id UUID PRIMARY KEY,
    edition edition_type NOT NULL,
    version VARCHAR(180) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);