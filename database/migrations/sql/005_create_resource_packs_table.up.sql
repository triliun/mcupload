-- Create resource_packs table
CREATE TABLE IF NOT EXISTS resource_packs (
    id UUID PRIMARY KEY,
    slug VARCHAR(180) NOT NULL UNIQUE,
    title VARCHAR(180) NOT NULL,
    content TEXT,
    status status_type NOT NULL DEFAULT 'draft',
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ
);

-- Create indexes for resource_packs
CREATE INDEX IF NOT EXISTS idx_resource_packs_author_id ON resource_packs(author_id);
CREATE INDEX IF NOT EXISTS idx_resource_packs_slug ON resource_packs(slug);