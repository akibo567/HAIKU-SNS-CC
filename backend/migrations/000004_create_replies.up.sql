CREATE TABLE replies (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id    UUID         NOT NULL REFERENCES haiku_posts(id) ON DELETE CASCADE,
    user_id    UUID         NOT NULL REFERENCES users(id)       ON DELETE CASCADE,
    ku1        VARCHAR(100) NOT NULL,
    ku2        VARCHAR(100) NOT NULL,
    ku3        VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_replies_post_id ON replies(post_id, created_at ASC);
