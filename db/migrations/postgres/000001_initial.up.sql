BEGIN;

CREATE TABLE person (
    id SERIAL UNIQUE PRIMARY KEY,
    twitch_id bigint UNIQUE,
    twitch_username text,
    telegram_id bigint UNIQUE,
    telegram_username text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone
);

CREATE TABLE token (
    twitch_id bigint UNIQUE,
    twitch_auth_code text,
    twitch_bearer text,
    twitch_bearer_expires_at timestamp with time zone,
    twitch_refresh_token text,
    invite_key text
);

CREATE TABLE settings (
    id SERIAL UNIQUE PRIMARY KEY,
    chat_telegram_id bigint UNIQUE NOT NULL,
    is_downloader_enabled boolean DEFAULT true
);

CREATE TABLE chat (
    twitch_id bigint UNIQUE,
    chat_telegram_id bigint NOT NULL,
    disabled boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone
);

CREATE TABLE invite (
    twitch_id bigint NOT NULL,
    chat_telegram_id bigint NOT NULL,
    chat_telegram_invite_link text,

    CONSTRAINT twitch_id_and_chat_id_unique UNIQUE (twitch_id, chat_telegram_id)
);

COMMIT;