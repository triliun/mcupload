-- Create pack_categories table
CREATE TABLE IF NOT EXISTS pack_categories (
    id UUID PRIMARY KEY,
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    pack_id UUID NOT NULL REFERENCES resource_packs(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add unique constraint
ALTER TABLE pack_categories ADD CONSTRAINT pack_id_category_id_unique UNIQUE (pack_id, category_id);

-- Create index
CREATE INDEX IF NOT EXISTS idx_pack_categories_pack_id ON pack_categories(author_id, pack_id);