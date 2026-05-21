-- 004_market_schema.sql
-- 广场集市模块：商品表

CREATE TABLE items (
    id           BIGINT PRIMARY KEY,
    seller_id    BIGINT       NOT NULL REFERENCES users(id),
    community_id BIGINT       NOT NULL REFERENCES communities(id),
    title        VARCHAR(100) NOT NULL,
    price        INT          NOT NULL,
    condition    VARCHAR(20)  NOT NULL DEFAULT 'like_new',
    category     VARCHAR(20)  NOT NULL DEFAULT 'other',
    images       TEXT[]       NOT NULL DEFAULT '{}',
    status       VARCHAR(20)  NOT NULL DEFAULT 'selling',
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_items_community_id ON items (community_id);
CREATE INDEX idx_items_community_category ON items (community_id, category);
CREATE INDEX idx_items_seller_id ON items (seller_id);
CREATE INDEX idx_items_created_at ON items (created_at DESC);
