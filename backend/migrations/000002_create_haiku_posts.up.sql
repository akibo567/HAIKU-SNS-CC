CREATE TABLE IF NOT EXISTS haiku_posts (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ku1        TEXT        NOT NULL,
    ku2        TEXT        NOT NULL,
    ku3        TEXT        NOT NULL,
    like_count INTEGER     NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_haiku_posts_user_id    ON haiku_posts(user_id);
CREATE INDEX IF NOT EXISTS idx_haiku_posts_created_at ON haiku_posts(created_at DESC);

CREATE TABLE IF NOT EXISTS likes (
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id    UUID        NOT NULL REFERENCES haiku_posts(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, post_id)
);
