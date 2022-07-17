CREATE TYPE transaction_type as ENUM (
    'income',
    'expense'
);

CREATE TYPE currency_type as ENUM (
    'TRY',
    'USD',
    'EUR'
);

CREATE TABLE transactions (
    transaction_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID  NOT NULL REFERENCES users,
    account_id UUID NOT NULL REFERENCES accounts,
    category_id UUID NOT NULL REFERENCES categories,

    date TIMESTAMP NOT NULL,
    type transaction_type NOT NULL,
    currency currency_type NOT NULL,
    amount INTEGER NOT NULL,
    notes TEXT NOT NULL DEFAULT '',

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);