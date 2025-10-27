-- Create enum types
CREATE TYPE role_type AS ENUM ('user', 'admin', 'ceo');
CREATE TYPE edition_type AS ENUM ('Java', 'Bedrock');
CREATE TYPE status_type AS ENUM ('draft', 'published', 'archived');