CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    login         VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (login, password_hash) VALUES ('admin', '$2a$10$e.Q6kFnSA591Gxi4tfx/LuyS7.NjEpFRLDvrnmuqHNILxgfHOpdvi');

CREATE TABLE IF NOT EXISTS secrets (
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title         VARCHAR(255) NOT NULL,
    metadata      TEXT DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    secret_type   VARCHAR(20) NOT NULL CHECK (secret_type IN ('login_password', 'text', 'binary', 'card')),

    -- login_password
    login         VARCHAR(255) DEFAULT '', -- для login_password
    password_hash TEXT DEFAULT '',         -- хэш/зашифрованный пароль

    --text
    text_data     TEXT DEFAULT '',         -- произвольный текст

    --binary
    binary_data   BYTEA DEFAULT '\x',      -- произвольные бинарные данные

    --card
    card_number   TEXT DEFAULT '',         -- номер карты
    card_exp      VARCHAR(10) DEFAULT '',  -- срок действия
    card_holder   TEXT DEFAULT '',          -- владелец карты

    CONSTRAINT secrets_user_title_uniq UNIQUE (user_id, title)
);
