BEGIN;

CREATE TABLE person (
    twitch_id bigint PRIMARY KEY,
    twitch_username text,
    telegram_id bigint UNIQUE,
    telegram_username text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone
);

CREATE TABLE token (
    twitch_id bigint PRIMARY KEY,
    twitch_auth_code text,
    twitch_bearer text,
    twitch_bearer_expires_at timestamp with time zone,
    twitch_refresh_token text,
    invite_key text
);

CREATE TABLE subchat (
    twitch_id bigint NOT NULL,
    subchat_telegram_id bigint NOT NULL,
    disabled boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone
);

CREATE TABLE invite (
    twitch_id bigint NOT NULL,
    subchat_telegram_id bigint NOT NULL,
    subchat_telegram_invite_link text
);

ALTER TABLE
    invite
ADD
    CONSTRAINT twitch_id_and_subchat_id_unique UNIQUE (twitch_id, subchat_telegram_id);

COMMIT;