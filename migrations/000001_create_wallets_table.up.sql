CREATE TABLE wallets (
    id TEXT PRIMARY KEY    CHECK (id != ''),
    balance FLOAT NOT NULL CHECK (balance >= 0)
);

