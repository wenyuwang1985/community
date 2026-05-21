-- 003_posts_schema.sql
-- 广场动态模块：帖子、评论、点赞

-- 帖子表
CREATE TABLE posts (
    id            BIGINT PRIMARY KEY,
    user_id       BIGINT       NOT NULL REFERENCES users(id),
    community_id  BIGINT       NOT NULL REFERENCES communities(id),
    tag           VARCHAR(20)  NOT NULL DEFAULT 'share',
    content       TEXT         NOT NULL,
    images        TEXT[]       NOT NULL DEFAULT '{}',
    status        VARCHAR(20)  NOT NULL DEFAULT 'normal',
    like_count    INT          NOT NULL DEFAULT 0,
    comment_count INT          NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_posts_community_id ON posts (community_id);
CREATE INDEX idx_posts_user_id ON posts (user_id);
CREATE INDEX idx_posts_created_at ON posts (created_at DESC);

-- 评论表
CREATE TABLE comments (
    id         BIGINT PRIMARY KEY,
    post_id    BIGINT       NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id    BIGINT       NOT NULL REFERENCES users(id),
    content    TEXT         NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_comments_post_id ON comments (post_id);
CREATE INDEX idx_comments_created_at ON comments (created_at);

-- 点赞表
CREATE TABLE post_likes (
    post_id    BIGINT      NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id    BIGINT      NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (post_id, user_id)
);

CREATE INDEX idx_post_likes_user_id ON post_likes (user_id);
