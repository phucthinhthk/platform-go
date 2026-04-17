-- Tự động chạy lần đầu khi container chưa có volume db_data

CREATE DATABASE IF NOT EXISTS platform_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE platform_db;

-- CREATE TABLE IF NOT EXISTS users (
--   id INT AUTO_INCREMENT PRIMARY KEY,
--   username VARCHAR(50) NOT NULL UNIQUE,
--   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

-- INSERT INTO users (username) VALUES ('admin');
