CREATE DATABASE airplanesco;

CREATE TABLE IF NOT EXISTS Customers (
    id SERIAL PRIMARY KEY,
    login VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50),
    surname VARCHAR(100),
    name VARCHAR(100)
);