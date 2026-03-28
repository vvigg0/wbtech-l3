CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    text TEXT NOT NULL,
    tg_id BIGINT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    send_at TIMESTAMP NOT NULL
);