ALTER TABLE users
    DROP INDEX users_email_unique,
    DROP COLUMN active,
    DROP COLUMN account_type,
    DROP COLUMN password_hash,
    DROP COLUMN email;
