-- Create pack_files table
CREATE TABLE IF NOT EXISTS pack_files (
    id UUID PRIMARY KEY,
    author_id UUID NOT NULL REFERENCES users(id),
    pack_id UUID NOT NULL REFERENCES resource_packs(id),
    filename VARCHAR(180) NOT NULL,
    size BIGINT NOT NULL,
    extension VARCHAR(10) NOT NULL,
    path TEXT NOT NULL UNIQUE,
    version_file VARCHAR(50) NOT NULL,
    minecraft_version_id UUID NOT NULL REFERENCES minecraft_versions(id),
    hash VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add unique constraint
ALTER TABLE pack_files ADD CONSTRAINT path_unique UNIQUE (path);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_pack_files_pack_id ON pack_files(pack_id);
CREATE INDEX IF NOT EXISTS idx_pack_files_minecraft_version ON pack_files(minecraft_version_id);