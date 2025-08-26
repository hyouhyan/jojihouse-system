ALTER TABLE users RENAME COLUMN contact TO discord_id;

ALTER TABLE users ALTER COLUMN discord_id TYPE BIGINT USING discord_id::BIGINT;

-- ALTER TABLE users ALTER COLUMN discord_id SET UNIQUE;
ALTER TABLE users ADD CONSTRAINT unique_discord_id UNIQUE (discord_id);
