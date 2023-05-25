CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    mnemonic VARCHAR(255) NOT NULL,
    xpriv VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL
);