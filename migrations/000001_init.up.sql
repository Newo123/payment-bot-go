CREATE SCHEMA bot;

CREATE TYPE bot.card_status AS ENUM('ENABLED', 'DISABLED');
CREATE TYPE bot.type_enum AS ENUM('DEPOSIT', 'PAYOUT');

CREATE TABLE bot.cards (
    id UUID PRIMARY KEY,
    number BIGINT NOT NULL UNIQUE,
    holder VARCHAR(1000) NOT NULL,
    status bot.card_status DEFAULT 'ENABLED',
    card_limit BIGINT NOT NULL,
    updated_at TIMESTAMP DEFAULT now() NOT NULL,
    created_at TIMESTAMP DEFAULT now() NOT NULL
);

CREATE TABLE bot.users (
    id UUID PRIMARY KEY,
    phone VARCHAR(255) NOT NULL UNIQUE,
    telegram_id BIGINT NOT NULL UNIQUE,
    language_code VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    username VARCHAR(255),
    updated_at TIMESTAMP DEFAULT now() NOT NULL,
    created_at TIMESTAMP DEFAULT now() NOT NULL
);

CREATE TABLE bot.transactions (
    id UUID PRIMARY KEY,
    amount BIGINT NOT NULL,
    user_id UUID NOT NULL,
    type bot.type_enum NOT NULL,
    operation_id BIGINT,
    card_id UUID NOT NULL,
    updated_at TIMESTAMP DEFAULT now() NOT NULL,
    created_at TIMESTAMP DEFAULT now() NOT NULL,
    CONSTRAINT fk_transactions_user FOREIGN KEY (user_id) REFERENCES bot.users(id),
    CONSTRAINT fk_transactions_card FOREIGN KEY (card_id) REFERENCES bot.cards(id)
);

CREATE INDEX idx_transactions_user_id ON bot.transactions(user_id);
CREATE INDEX idx_transactions_card_id ON bot.transactions(card_id);
CREATE INDEX idx_transactions_created_at ON bot.transactions(created_at);
CREATE INDEX idx_cards_status ON bot.cards(status);
CREATE INDEX idx_users_telegram_id ON bot.users(telegram_id);