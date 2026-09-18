ALTER TABLE users
    ADD COLUMN email VARCHAR(255) NULL AFTER name,
    ADD COLUMN password_hash VARCHAR(255) NULL AFTER email,
    ADD COLUMN account_type VARCHAR(50) NOT NULL DEFAULT 'admin' AFTER password_hash,
    ADD COLUMN active BOOLEAN NOT NULL DEFAULT TRUE AFTER account_type,
    ADD UNIQUE KEY users_email_unique (email);
