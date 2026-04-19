-- +goose Up
-- Agents table
CREATE TABLE agents (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    model TEXT NOT NULL,
    provider TEXT NOT NULL,
    persona TEXT,
    api_key_ref TEXT, -- Reference to secret manager or encrypted store
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Rooms (Conversations) table
CREATE TABLE rooms (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    goal TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Messages table
CREATE TABLE messages (
    id TEXT PRIMARY KEY,
    room_id TEXT NOT NULL,
    agent_id TEXT, -- Can be NULL for user messages
    role TEXT NOT NULL, -- 'user', 'assistant', 'system'
    content TEXT NOT NULL,
    tokens_used INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
    FOREIGN KEY (agent_id) REFERENCES agents(id) ON DELETE SET NULL
);
CREATE INDEX idx_messages_room_id ON messages(room_id);
CREATE INDEX idx_messages_room_id_created_at ON messages(room_id, created_at);
CREATE INDEX idx_messages_agent_id ON messages(agent_id);

-- Usage Log table
CREATE TABLE usage_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    provider TEXT NOT NULL,
    model TEXT NOT NULL,
    prompt_tokens INTEGER NOT NULL,
    completion_tokens INTEGER NOT NULL,
    total_tokens INTEGER NOT NULL,
    latency_ms INTEGER NOT NULL,
    cost REAL DEFAULT 0.0,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_usage_log_timestamp ON usage_log(timestamp);

-- +goose Down
DROP TABLE usage_log;
DROP TABLE messages;
DROP TABLE rooms;
DROP TABLE agents;
