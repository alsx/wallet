CREATE DATABASE wallet;

\connect wallet

CREATE TABLE IF NOT EXISTS accounts
(
    id      varchar(40) NOT NULL,
    balance money CHECK (balance::numeric::float8 > 0),
    PRIMARY KEY(id),
    UNIQUE(id)
);
CREATE TABLE IF NOT EXISTS payments
(
    from_account varchar(40) REFERENCES accounts(id),
    amount       money CHECK (amount::numeric::float8 > 0),
    to_account   varchar(40) REFERENCES accounts(id)
);

INSERT INTO accounts VALUES ('bob123', 100.00), ('alice456', 0.01);
INSERT INTO payments VALUES ('alice456', 50.00, 'bob123');
