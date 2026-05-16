-- 001_init_schema.sql
-- 核心基础表：communities, users, user_community_subscriptions

CREATE EXTENSION IF NOT EXISTS postgis;

-- 街镇表
CREATE TABLE communities (
    id         BIGINT PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    province   VARCHAR(50)  NOT NULL,
    city       VARCHAR(50)  NOT NULL,
    district   VARCHAR(50)  NOT NULL,
    adcode     VARCHAR(12),
    boundary   GEOMETRY(MultiPolygon, 4326),
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_communities_adcode ON communities (adcode);

-- 用户表
CREATE TABLE users (
    id           BIGINT PRIMARY KEY,
    phone        VARCHAR(20)  NOT NULL UNIQUE,
    nickname     VARCHAR(50)  NOT NULL DEFAULT '',
    avatar_url   VARCHAR(500) NOT NULL DEFAULT '',
    credit_score SMALLINT     NOT NULL DEFAULT 100,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now()
);

-- 用户订阅关系表
CREATE TABLE user_community_subscriptions (
    id           BIGINT PRIMARY KEY,
    user_id      BIGINT      NOT NULL REFERENCES users(id),
    community_id BIGINT      NOT NULL REFERENCES communities(id),
    is_primary   BOOLEAN     NOT NULL DEFAULT false,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id, community_id)
);

CREATE INDEX idx_subscriptions_user_id ON user_community_subscriptions (user_id);
CREATE INDEX idx_subscriptions_community_id ON user_community_subscriptions (community_id);
