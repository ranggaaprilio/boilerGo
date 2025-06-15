-- Create database if not exists
CREATE DATABASE IF NOT EXISTS boilergo;

-- Use the database
USE boilergo;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Grant privileges to the user
GRANT ALL PRIVILEGES ON boilergo.* TO 'rangga'@'%';
FLUSH PRIVILEGES;
