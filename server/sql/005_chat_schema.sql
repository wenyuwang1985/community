-- 005_chat_schema.sql
-- 聊天模块：会话、参与者、消息

-- 会话表
CREATE TABLE conversations (
    id           BIGINT PRIMARY KEY,
    type         VARCHAR(20)  NOT NULL CHECK (type IN ('private', 'channel')),
    community_id BIGINT       REFERENCES communities(id),
    name         VARCHAR(100) NOT NULL DEFAULT '',
    created_by   BIGINT       REFERENCES users(id),
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_conversations_community_id ON conversations (community_id);

-- 会话参与者表
CREATE TABLE conversation_participants (
    conversation_id BIGINT      NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    user_id         BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (conversation_id, user_id)
);

CREATE INDEX idx_participants_user_id ON conversation_participants (user_id);
CREATE INDEX idx_participants_conversation_id ON conversation_participants (conversation_id);

-- 消息表
CREATE TABLE messages (
    id              BIGINT PRIMARY KEY,
    conversation_id BIGINT       NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id       BIGINT       NOT NULL REFERENCES users(id),
    content         TEXT         NOT NULL,
    type            VARCHAR(20)  NOT NULL DEFAULT 'text' CHECK (type IN ('text', 'image', 'system')),
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_messages_conversation_id ON messages (conversation_id);
CREATE INDEX idx_messages_created_at ON messages (created_at);
