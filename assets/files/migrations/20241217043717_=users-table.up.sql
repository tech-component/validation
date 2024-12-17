-- Enable the uuid-ossp extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT NOT NULL,
    password TEXT NOT NULL
);

ALTER TABLE users ADD CONSTRAINT users_unique_email UNIQUE(email);